package cmd

import (
	"context"

	"github.com/loft-sh/devpod/pkg/log"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/options"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/upcloud"
	"github.com/spf13/cobra"
)

// StopCmd holds the cmd flags
type StopCmd struct{}

// NewStopCmd defines a command
func NewStopCmd() *cobra.Command {
	cmd := &StopCmd{}
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv(false)
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return stopCmd
}

// Run runs the command logic
func (cmd *StopCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	client := upcloud.NewUpCloud(options.Username, options.Password)
	
	log.Infof("Stopping server %s...", options.MachineID)
	err := client.Stop(ctx, options.MachineID)
	if err != nil {
		return err
	}
	
	log.Infof("Successfully stopped server %s", options.MachineID)
	return nil
}