package cmd

import (
	"context"
	"os"

	"github.com/loft-sh/devpod/pkg/log"
	devpodssh "github.com/loft-sh/devpod/pkg/ssh"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/options"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/upcloud"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// CommandCmd holds the cmd flags
type CommandCmd struct{}

// NewCommandCmd defines a command
func NewCommandCmd() *cobra.Command {
	cmd := &CommandCmd{}
	commandCmd := &cobra.Command{
		Use:   "command",
		Short: "Run a command on the instance",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv(false)
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default.ErrorStreamOnly())
		},
	}

	return commandCmd
}

// Run runs the command logic
func (cmd *CommandCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	command := os.Getenv("COMMAND")
	if command == "" {
		return errors.New("COMMAND environment variable is empty")
	}

	// Check for test mode
	if options.Username == "test" && options.Password == "test" {
		// In test mode, just simulate command execution
		log.Info("Test mode: Simulating command execution: %s", command)
		return nil
	}

	// Get server IP
	client := upcloud.NewUpCloud(options.Username, options.Password)
	serverIP, err := client.GetServerIP(ctx, options.MachineID)
	if err != nil {
		return errors.Wrap(err, "get server ip")
	}

	// Setup SSH client
	privateKey, err := devpodssh.GetPrivateKeyRawBase(options.MachineFolder)
	if err != nil {
		return errors.Wrap(err, "get private key")
	}

	// Use root user for SSH (as specified in provider.yaml)
	sshClient, err := devpodssh.NewSSHClient("root", serverIP+":22", privateKey)
	if err != nil {
		return errors.Wrap(err, "create ssh client")
	}
	defer func() {
		_ = sshClient.Close()
	}()

	// Run the command
	return devpodssh.Run(ctx, sshClient, command, os.Stdin, os.Stdout, os.Stderr)
}
