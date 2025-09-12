package cmd

import (
	"context"
	"encoding/base64"
	
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/options"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/upcloud"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// CreateCmd holds the cmd flags
type CreateCmd struct{}

// NewCreateCmd defines a command
func NewCreateCmd() *cobra.Command {
	cmd := &CreateCmd{}
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv(false)
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return createCmd
}

// Run runs the command logic
func (cmd *CreateCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	client := upcloud.NewUpCloud(options.Username, options.Password)
	
	// Get SSH public key
	publicKeyBase, err := ssh.GetPublicKeyBase(options.MachineFolder)
	if err != nil {
		return errors.Wrap(err, "get public key")
	}
	
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase)
	if err != nil {
		return errors.Wrap(err, "decode public key")
	}
	
	// Create server configuration
	serverConfig := &upcloud.ServerConfig{
		Hostname:  options.MachineID,
		Zone:      options.Zone,
		Plan:      options.Plan,
		Storage:   options.Storage,
		Image:     options.Image,
		SSHKey:    string(publicKey),
		UserData:  GetCloudInitScript(options.MachineID),
	}
	
	// Create the server
	log.Infof("Creating UpCloud server %s...", options.MachineID)
	err = client.Create(ctx, serverConfig)
	if err != nil {
		return errors.Wrap(err, "create server")
	}
	
	log.Infof("Successfully created server %s", options.MachineID)
	return nil
}

func GetCloudInitScript(machineID string) string {
	return `#!/bin/bash
set -e

# Create devpod user
useradd -m -s /bin/bash devpod || true
usermod -aG docker devpod 2>/dev/null || true
usermod -aG sudo devpod
echo "devpod ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/devpod

# Setup SSH directory
mkdir -p /home/devpod/.ssh
chmod 700 /home/devpod/.ssh
chown -R devpod:devpod /home/devpod

# Install Docker if not present
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com | sh
    systemctl enable docker
    systemctl start docker
fi

# Ensure required directories exist
mkdir -p /opt/devpod
chown -R devpod:devpod /opt/devpod
`
}