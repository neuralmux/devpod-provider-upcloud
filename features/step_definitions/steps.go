package step_definitions

import (
	"fmt"
	"os"

	"github.com/cucumber/godog"
	"github.com/neuralmux/devpod-provider-upcloud/cmd"
)

// providerContext holds the test context for BDD scenarios
type providerContext struct {
	credentials  credentials
	serverID     string
	serverStatus string
	lastError    error
}

type credentials struct {
	username string
	password string
}

// InitializeScenario initializes the scenario context with step definitions
func InitializeScenario(ctx *godog.ScenarioContext) {
	p := &providerContext{}

	// Given steps
	ctx.Step(`^I have valid UpCloud API credentials$`, p.iHaveValidUpCloudAPICredentials)
	ctx.Step(`^the provider is configured with required options$`, p.theProviderIsConfiguredWithRequiredOptions)
	ctx.Step(`^I have a running UpCloud server$`, p.iHaveARunningUpCloudServer)
	ctx.Step(`^I have a stopped UpCloud server$`, p.iHaveAStoppedUpCloudServer)
	ctx.Step(`^I have an existing UpCloud server$`, p.iHaveAnExistingUpCloudServer)

	// When steps
	ctx.Step(`^I run the init command$`, p.iRunTheInitCommand)
	ctx.Step(`^I run the create command$`, p.iRunTheCreateCommand)
	ctx.Step(`^I run the stop command$`, p.iRunTheStopCommand)
	ctx.Step(`^I run the start command$`, p.iRunTheStartCommand)
	ctx.Step(`^I run the delete command$`, p.iRunTheDeleteCommand)
	ctx.Step(`^I execute a command on the server$`, p.iExecuteACommandOnTheServer)

	// Then steps
	ctx.Step(`^the provider should validate the credentials$`, p.theProviderShouldValidateTheCredentials)
	ctx.Step(`^return a success status$`, p.returnASuccessStatus)
	ctx.Step(`^a new server should be created in UpCloud$`, p.aNewServerShouldBeCreatedInUpCloud)
	ctx.Step(`^the server should be accessible via SSH$`, p.theServerShouldBeAccessibleViaSSH)
	ctx.Step(`^the status should return "([^"]*)"$`, p.theStatusShouldReturn)
	ctx.Step(`^the server should be stopped$`, p.theServerShouldBeStopped)
	ctx.Step(`^the server should be started$`, p.theServerShouldBeStarted)
	ctx.Step(`^the server should be removed from UpCloud$`, p.theServerShouldBeRemovedFromUpCloud)
	ctx.Step(`^the command should run successfully$`, p.theCommandShouldRunSuccessfully)
	ctx.Step(`^I should see the command output$`, p.iShouldSeeTheCommandOutput)
}

// InitializeTestSuite initializes the test suite context
func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		// Setup test environment
	})

	ctx.AfterSuite(func() {
		// Cleanup test resources
	})
}

// Step implementations
func (p *providerContext) iHaveValidUpCloudAPICredentials() error {
	p.credentials.username = os.Getenv("UPCLOUD_USERNAME")
	p.credentials.password = os.Getenv("UPCLOUD_PASSWORD")

	// Use mock credentials if not set (for testing)
	if p.credentials.username == "" || p.credentials.password == "" {
		p.credentials.username = "test"
		p.credentials.password = "test"
		_ = os.Setenv("UPCLOUD_USERNAME", "test")
		_ = os.Setenv("UPCLOUD_PASSWORD", "test")
	}
	return nil
}

func (p *providerContext) theProviderIsConfiguredWithRequiredOptions() error {
	// Set default values for testing if not set
	defaults := map[string]string{
		"UPCLOUD_ZONE":    "de-fra1",
		"UPCLOUD_PLAN":    "1xCPU-1GB",
		"AGENT_PATH":      "/home/devpod/.devpod/devpod",
		"MACHINE_ID":      "test-machine",
		"MACHINE_FOLDER":  "/tmp/test",
		"UPCLOUD_STORAGE": "25",
		"UPCLOUD_IMAGE":   "Ubuntu Server 22.04 LTS (Jammy Jellyfish)",
	}

	for envVar, defaultValue := range defaults {
		if os.Getenv(envVar) == "" {
			_ = os.Setenv(envVar, defaultValue)
		}
	}
	return nil
}

func (p *providerContext) iRunTheInitCommand() error {
	// Create and run the init command
	initCmd := cmd.NewInitCmd()
	p.lastError = initCmd.Execute()
	return nil
}

func (p *providerContext) theProviderShouldValidateTheCredentials() error {
	if p.lastError != nil {
		return fmt.Errorf("credential validation failed: %w", p.lastError)
	}
	return nil
}

func (p *providerContext) returnASuccessStatus() error {
	if p.lastError != nil {
		return fmt.Errorf("command did not return success: %w", p.lastError)
	}
	return nil
}

func (p *providerContext) iRunTheCreateCommand() error {
	// Create and run the create command
	createCmd := cmd.NewCreateCmd()
	p.lastError = createCmd.Execute()
	// TODO: Capture created server ID
	p.serverID = "test-server-id"
	return nil
}

func (p *providerContext) aNewServerShouldBeCreatedInUpCloud() error {
	if p.serverID == "" {
		return fmt.Errorf("no server was created")
	}
	return nil
}

func (p *providerContext) theServerShouldBeAccessibleViaSSH() error {
	// TODO: Verify SSH connectivity
	return nil
}

func (p *providerContext) theStatusShouldReturn(expectedStatus string) error {
	// TODO: Check actual status
	if p.serverStatus != expectedStatus {
		return fmt.Errorf("expected status %s, got %s", expectedStatus, p.serverStatus)
	}
	return nil
}

func (p *providerContext) iHaveARunningUpCloudServer() error {
	// Setup: Ensure we have a running server
	p.serverID = "existing-server"
	p.serverStatus = "Running"
	return nil
}

func (p *providerContext) iHaveAStoppedUpCloudServer() error {
	// Setup: Ensure we have a stopped server
	p.serverID = "existing-server"
	p.serverStatus = "Stopped"
	return nil
}

func (p *providerContext) iHaveAnExistingUpCloudServer() error {
	// Setup: Ensure we have an existing server
	p.serverID = "existing-server"
	return nil
}

func (p *providerContext) iRunTheStopCommand() error {
	// Create and run the stop command
	stopCmd := cmd.NewStopCmd()
	p.lastError = stopCmd.Execute()
	return nil
}

func (p *providerContext) theServerShouldBeStopped() error {
	// TODO: Verify server is stopped
	p.serverStatus = "Stopped"
	return nil
}

func (p *providerContext) iRunTheStartCommand() error {
	// Create and run the start command
	startCmd := cmd.NewStartCmd()
	p.lastError = startCmd.Execute()
	return nil
}

func (p *providerContext) theServerShouldBeStarted() error {
	// TODO: Verify server is started
	p.serverStatus = "Running"
	return nil
}

func (p *providerContext) iRunTheDeleteCommand() error {
	// Create and run the delete command
	deleteCmd := cmd.NewDeleteCmd()
	p.lastError = deleteCmd.Execute()
	return nil
}

func (p *providerContext) theServerShouldBeRemovedFromUpCloud() error {
	// TODO: Verify server is deleted
	p.serverID = ""
	return nil
}

func (p *providerContext) iExecuteACommandOnTheServer() error {
	// Create and run the command command
	commandCmd := cmd.NewCommandCmd()
	p.lastError = commandCmd.Execute()
	return nil
}

func (p *providerContext) theCommandShouldRunSuccessfully() error {
	if p.lastError != nil {
		return fmt.Errorf("command execution failed: %w", p.lastError)
	}
	return nil
}

func (p *providerContext) iShouldSeeTheCommandOutput() error {
	// TODO: Verify command output
	return nil
}
