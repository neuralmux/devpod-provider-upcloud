package cmd

import (
	"context"
	"fmt"

	"github.com/loft-sh/devpod/pkg/log"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/options"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/upcloud"
	"github.com/spf13/cobra"
)

// StatusCmd holds the cmd flags
type StatusCmd struct{}

// NewStatusCmd defines a command
func NewStatusCmd() *cobra.Command {
	cmd := &StatusCmd{}
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Retrieve the status of an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv(false)
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return statusCmd
}

// Run runs the command logic
func (cmd *StatusCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	client := upcloud.NewUpCloud(options.Username, options.Password)

	status, err := client.Status(ctx, options.MachineID)
	if err != nil {
		fmt.Println("NotFound")
		return nil
	}

	// Print status for DevPod to consume
	fmt.Println(status)
	return nil
}
