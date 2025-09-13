package cmd

import (
	"context"

	"github.com/loft-sh/devpod/pkg/log"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/options"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/upcloud"
	"github.com/spf13/cobra"
)

// DeleteCmd holds the cmd flags
type DeleteCmd struct{}

// NewDeleteCmd defines a command
func NewDeleteCmd() *cobra.Command {
	cmd := &DeleteCmd{}
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv(false)
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return deleteCmd
}

// Run runs the command logic
func (cmd *DeleteCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	client := upcloud.NewUpCloud(options.Username, options.Password)

	log.Infof("Deleting server %s...", options.MachineID)
	err := client.Delete(ctx, options.MachineID)
	if err != nil {
		log.Warnf("Failed to delete server %s: %v", options.MachineID, err)
		// Don't fail if server doesn't exist
		return nil
	}

	log.Infof("Successfully deleted server %s", options.MachineID)
	return nil
}
