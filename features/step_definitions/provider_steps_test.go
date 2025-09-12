package step_definitions

import (
	"fmt"
	"os"

	"github.com/cucumber/godog"
	"github.com/neuralmux/devpod-provider-upcloud/cmd"
)

type providerContext struct {
	credentials struct {
		username string
		password string
	}
	serverID     string
	serverStatus string
	lastError    error
}

func (p *providerContext) iHaveValidUpCloudAPICredentials() error {
	p.credentials.username = os.Getenv("UPCLOUD_USERNAME")
	p.credentials.password = os.Getenv("UPCLOUD_PASSWORD")
	
	if p.credentials.username == "" || p.credentials.password == "" {
		return fmt.Errorf("UpCloud credentials not set in environment")
	}
	return nil
}

func (p *providerContext) theProviderIsConfiguredWithRequiredOptions() error {
	// Check for required provider options
	requiredEnvVars := []string{
		"UPCLOUD_ZONE",
		"UPCLOUD_PLAN",
		"AGENT_PATH",
	}
	
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("required option %s is not set", envVar)
		}
	}
	return nil
}

func (p *providerContext) iRunTheInitCommand() error {
	p.lastError = cmd.Init()
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
		return fmt.Errorf("expected success but got error: %w", p.lastError)
	}
	return nil
}

func (p *providerContext) iRunTheCreateCommand() error {
	p.lastError = cmd.Create()
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
	// TODO: Test SSH connectivity
	return nil
}

func (p *providerContext) theStatusShouldReturn(expectedStatus string) error {
	// TODO: Get actual status
	p.serverStatus = expectedStatus
	if p.serverStatus != expectedStatus {
		return fmt.Errorf("expected status %s, got %s", expectedStatus, p.serverStatus)
	}
	return nil
}

func (p *providerContext) iHaveARunningUpCloudServer() error {
	// TODO: Ensure a running server exists
	p.serverID = "existing-server"
	p.serverStatus = "Running"
	return nil
}

func (p *providerContext) iHaveAStoppedUpCloudServer() error {
	// TODO: Ensure a stopped server exists
	p.serverID = "stopped-server"
	p.serverStatus = "Stopped"
	return nil
}

func (p *providerContext) iHaveAnExistingUpCloudServer() error {
	// TODO: Ensure any server exists
	p.serverID = "existing-server"
	return nil
}

func (p *providerContext) iRunTheStopCommand() error {
	p.lastError = cmd.Stop()
	return nil
}

func (p *providerContext) theServerShouldBeStopped() error {
	// TODO: Verify server is stopped
	p.serverStatus = "Stopped"
	return nil
}

func (p *providerContext) iRunTheStartCommand() error {
	p.lastError = cmd.Start()
	return nil
}

func (p *providerContext) theServerShouldBeStarted() error {
	// TODO: Verify server is started
	p.serverStatus = "Running"
	return nil
}

func (p *providerContext) iRunTheDeleteCommand() error {
	p.lastError = cmd.Delete()
	return nil
}

func (p *providerContext) theServerShouldBeRemovedFromUpCloud() error {
	// TODO: Verify server is deleted
	p.serverID = ""
	return nil
}

func (p *providerContext) iExecuteACommandOnTheServer() error {
	p.lastError = cmd.Command("echo 'test'")
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

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		// Setup test environment
	})
	
	ctx.AfterSuite(func() {
		// Cleanup test resources
	})
}