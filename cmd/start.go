package cmd

import (
	"context"

	"github.com/loft-sh/devpod/pkg/log"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/options"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/upcloud"
	"github.com/spf13/cobra"
)

// StartCmd holds the cmd flags
type StartCmd struct{}

// NewStartCmd defines a command
func NewStartCmd() *cobra.Command {
	cmd := &StartCmd{}
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv(false)
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return startCmd
}

// Run runs the command logic
func (cmd *StartCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	client := upcloud.NewUpCloud(options.Username, options.Password)
	
	log.Infof("Starting server %s...", options.MachineID)
	err := client.Start(ctx, options.MachineID)
	if err != nil {
		return err
	}
	
	log.Infof("Successfully started server %s", options.MachineID)
	return nil
}