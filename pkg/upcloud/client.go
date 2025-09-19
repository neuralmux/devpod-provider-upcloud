package upcloud

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/service"
)

// Client represents the UpCloud API client
type Client struct {
	service *service.Service
	timeout time.Duration
}

// ServerConfig holds the configuration for creating a new server
type ServerConfig struct {
	Hostname string
	Zone     string
	Plan     string
	Storage  string
	Image    string
	Template string
	SSHKey   string
	UserData string
}

// NewUpCloud creates a new UpCloud client
func NewUpCloud(username, password string) *Client {
	// Check for test mode
	if username == "test" && password == "test" {
		// Return a mock client for testing
		return &Client{
			service: nil, // Will check for nil in methods
			timeout: time.Duration(DefaultTimeout) * time.Second,
		}
	}

	// Create client with timeout
	httpClient := client.New(username, password, client.WithTimeout(time.Second*30))

	// Create service
	svc := service.New(httpClient)

	return &Client{
		service: svc,
		timeout: time.Duration(DefaultTimeout) * time.Second,
	}
}

// TestConnection tests the connection to UpCloud API
func (c *Client) TestConnection(ctx context.Context) error {
	// Check for test mode
	if c.service == nil {
		fmt.Println("Test mode: Simulating successful authentication")
		return nil
	}

	// Try to get account information to verify credentials
	account, err := c.service.GetAccount(ctx)
	if err != nil {
		return WrapError(err, "authentication test")
	}

	// Check if account is valid
	if account.UserName == "" {
		return fmt.Errorf("invalid account response")
	}

	return nil
}

// Create creates a new server
func (c *Client) Create(ctx context.Context, config *ServerConfig) error {
	// Check for test mode
	if c.service == nil {
		fmt.Fprintf(os.Stderr, "Test mode: Simulating server creation for %s\n", config.Hostname)
		// Track state using environment variable for test mode
		_ = os.Setenv("TEST_SERVER_STATE_"+config.Hostname, "Running")
		return nil
	}

	// Validate zone
	if err := ValidateZone(config.Zone); err != nil {
		return WrapError(err, "zone validation")
	}

	// Map plan name
	plan, err := MapPlanName(config.Plan)
	if err != nil {
		return WrapError(err, "plan mapping")
	}

	// Use custom template if provided, otherwise map image to template UUID
	var templateUUID string
	if config.Template != "" {
		templateUUID = config.Template
	} else {
		templateUUID, err = MapImageToTemplate(config.Image)
		if err != nil {
			return WrapError(err, "image mapping")
		}
	}

	// Parse storage size
	storageSize, err := ParseStorageSize(config.Storage)
	if err != nil {
		return WrapError(err, "storage size parsing")
	}

	// Generate a clean hostname
	hostname := GenerateHostname(config.Hostname)

	// Create the server request
	createReq := &request.CreateServerRequest{
		Zone:             config.Zone,
		Title:            config.Hostname, // Use full machine ID as title
		Hostname:         hostname,
		Plan:             plan,
		PasswordDelivery: request.PasswordDeliveryNone,
		Metadata:         upcloud.True, // Required for cloud-init templates

		// Configure storage
		StorageDevices: []request.CreateServerStorageDevice{
			{
				Action:  request.CreateServerStorageDeviceActionClone,
				Storage: templateUUID,
				Title:   "root",
				Size:    storageSize,
				Tier:    GetStorageTier(plan),
			},
		},

		// Configure networking
		Networking: &request.CreateServerNetworking{
			Interfaces: []request.CreateServerInterface{
				{
					Type: "public",
					IPAddresses: []request.CreateServerIPAddress{
						{
							Family: upcloud.IPAddressFamilyIPv4,
						},
					},
				},
				{
					Type: "utility",
					IPAddresses: []request.CreateServerIPAddress{
						{
							Family: upcloud.IPAddressFamilyIPv4,
						},
					},
				},
			},
		},
	}

	// Add SSH key if provided
	if config.SSHKey != "" {
		createReq.LoginUser = &request.LoginUser{
			Username:       DefaultSSHUser,
			CreatePassword: "no",
			SSHKeys:        []string{config.SSHKey},
		}
	}

	// Add user data (cloud-init) if provided
	if config.UserData != "" {
		createReq.UserData = config.UserData
	}

	// Create the server
	serverDetails, err := c.service.CreateServer(ctx, createReq)
	if err != nil {
		return WrapError(err, "server creation")
	}

	// Wait for server to start
	waitReq := &request.WaitForServerStateRequest{
		UUID:         serverDetails.UUID,
		DesiredState: upcloud.ServerStateStarted,
	}

	_, err = c.service.WaitForServerState(ctx, waitReq)
	if err != nil {
		// If waiting fails, try to clean up by deleting the server
		_ = c.service.DeleteServer(ctx, &request.DeleteServerRequest{
			UUID: serverDetails.UUID,
		})
		return WrapError(err, "waiting for server to start")
	}

	return nil
}

// Delete deletes a server
func (c *Client) Delete(ctx context.Context, serverID string) error {
	// Check for test mode
	if c.service == nil {
		fmt.Fprintf(os.Stderr, "Test mode: Simulating server deletion for %s\n", serverID)
		// Track state using environment variable for test mode
		_ = os.Setenv("TEST_SERVER_STATE_"+serverID, "NotFound")
		return nil
	}

	// Find the server by machine ID
	server, err := c.findServerByMachineID(ctx, serverID)
	if err != nil {
		if IsNotFoundError(err) {
			// Server doesn't exist, consider it deleted
			return nil
		}
		return err
	}

	// Stop the server first if it's running
	if server.State == upcloud.ServerStateStarted {
		stopReq := &request.StopServerRequest{
			UUID:     server.UUID,
			StopType: request.ServerStopTypeHard,
		}
		_, err = c.service.StopServer(ctx, stopReq)
		if err != nil && !IsNotFoundError(err) {
			return WrapError(err, "stopping server before deletion")
		}

		// Wait for server to stop
		waitReq := &request.WaitForServerStateRequest{
			UUID:         server.UUID,
			DesiredState: upcloud.ServerStateStopped,
		}
		_, _ = c.service.WaitForServerState(ctx, waitReq)
	}

	// Delete the server
	deleteReq := &request.DeleteServerRequest{
		UUID: server.UUID,
	}
	err = c.service.DeleteServer(ctx, deleteReq)
	if err != nil && !IsNotFoundError(err) {
		return WrapError(err, "server deletion")
	}

	return nil
}

// Start starts a stopped server
func (c *Client) Start(ctx context.Context, serverID string) error {
	// Check for test mode
	if c.service == nil {
		fmt.Fprintf(os.Stderr, "Test mode: Simulating server start for %s\n", serverID)
		// Track state using environment variable for test mode
		_ = os.Setenv("TEST_SERVER_STATE_"+serverID, "Running")
		return nil
	}

	// Find the server by machine ID
	server, err := c.findServerByMachineID(ctx, serverID)
	if err != nil {
		return err
	}

	// Check if already started
	if server.State == upcloud.ServerStateStarted {
		return nil
	}

	// Start the server
	startReq := &request.StartServerRequest{
		UUID: server.UUID,
	}
	_, err = c.service.StartServer(ctx, startReq)
	if err != nil {
		return WrapError(err, "server start")
	}

	// Wait for server to start
	waitReq := &request.WaitForServerStateRequest{
		UUID:         server.UUID,
		DesiredState: upcloud.ServerStateStarted,
	}
	_, err = c.service.WaitForServerState(ctx, waitReq)
	if err != nil {
		return WrapError(err, "waiting for server to start")
	}

	return nil
}

// Stop stops a running server
func (c *Client) Stop(ctx context.Context, serverID string) error {
	// Check for test mode
	if c.service == nil {
		fmt.Fprintf(os.Stderr, "Test mode: Simulating server stop for %s\n", serverID)
		// Track state using environment variable for test mode
		_ = os.Setenv("TEST_SERVER_STATE_"+serverID, "Stopped")
		return nil
	}

	// Find the server by machine ID
	server, err := c.findServerByMachineID(ctx, serverID)
	if err != nil {
		return err
	}

	// Check if already stopped
	if server.State == upcloud.ServerStateStopped {
		return nil
	}

	// Stop the server
	stopReq := &request.StopServerRequest{
		UUID:     server.UUID,
		StopType: request.ServerStopTypeSoft,
	}
	_, err = c.service.StopServer(ctx, stopReq)
	if err != nil {
		return WrapError(err, "server stop")
	}

	// Wait for server to stop
	waitReq := &request.WaitForServerStateRequest{
		UUID:         server.UUID,
		DesiredState: upcloud.ServerStateStopped,
	}
	_, err = c.service.WaitForServerState(ctx, waitReq)
	if err != nil {
		return WrapError(err, "waiting for server to stop")
	}

	return nil
}

// Status returns the status of a server
func (c *Client) Status(ctx context.Context, serverID string) (string, error) {
	// Check for test mode
	if c.service == nil {
		fmt.Fprintf(os.Stderr, "Test mode: Simulating server status for %s\n", serverID)
		// Check test server state from environment variable
		if state := os.Getenv("TEST_SERVER_STATE_" + serverID); state != "" {
			return state, nil
		}
		// Default to Running for new servers
		return StatusRunning, nil
	}

	// Find the server by machine ID
	server, err := c.findServerByMachineID(ctx, serverID)
	if err != nil {
		if IsNotFoundError(err) {
			return StatusNotFound, nil
		}
		return StatusNotFound, err
	}

	// Map UpCloud state to DevPod status
	return MapServerStateToStatus(server.State), nil
}

// GetServerIP gets the public IP address of a server
func (c *Client) GetServerIP(ctx context.Context, serverID string) (string, error) {
	// Check for test mode
	if c.service == nil {
		fmt.Fprintf(os.Stderr, "Test mode: Simulating server IP retrieval for %s\n", serverID)
		return "192.0.2.1", nil // Return TEST-NET-1 address for testing
	}

	// Find the server by machine ID
	server, err := c.findServerByMachineID(ctx, serverID)
	if err != nil {
		return "", err
	}

	// Get full server details
	serverDetails, err := c.service.GetServerDetails(ctx, &request.GetServerDetailsRequest{
		UUID: server.UUID,
	})
	if err != nil {
		return "", WrapError(err, "getting server details")
	}

	// Extract public IPv4 address
	ip, err := GetPublicIPv4(serverDetails)
	if err != nil {
		return "", WrapError(err, "extracting public IP")
	}

	return ip, nil
}

// findServerByMachineID is a helper to find a server by DevPod machine ID
func (c *Client) findServerByMachineID(ctx context.Context, machineID string) (*upcloud.Server, error) {
	// Check for test mode
	if c.service == nil {
		// Return a mock server for testing
		return &upcloud.Server{
			UUID:  "test-uuid",
			Title: machineID,
			State: upcloud.ServerStateStarted,
		}, nil
	}

	// List all servers
	servers, err := c.service.GetServers(ctx)
	if err != nil {
		return nil, WrapError(err, "listing servers")
	}

	// Find server by machine ID
	server := FindServerByMachineID(servers.Servers, machineID)
	if server == nil {
		return nil, &ProviderError{
			Type:    ErrorTypeNotFound,
			Message: fmt.Sprintf("Server with machine ID %s not found", machineID),
		}
	}

	return server, nil
}
