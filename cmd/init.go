package cmd

import (
	"context"
	"fmt"

	"github.com/loft-sh/devpod/pkg/log"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/options"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/upcloud"
	"github.com/spf13/cobra"
)

// InitCmd holds the cmd flags
type InitCmd struct{}

// NewInitCmd defines a init command
func NewInitCmd() *cobra.Command {
	cmd := &InitCmd{}
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Init a machine",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv(true)
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return initCmd
}

// Run runs the command logic
func (cmd *InitCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	// Validate UpCloud credentials
	if options.Username == "" || options.Password == "" {
		return fmt.Errorf("UPCLOUD_USERNAME and UPCLOUD_PASSWORD must be set")
	}

	// Test API connection
	client := upcloud.NewUpCloud(options.Username, options.Password)
	err := client.TestConnection(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to UpCloud API: %w", err)
	}

	log.Infof("Successfully initialized UpCloud provider")
	return nil
}
