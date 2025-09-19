package options

import (
	"fmt"
	"os"
)

type Options struct {
	MachineID     string
	MachineFolder string

	Zone     string
	Plan     string
	Storage  string
	Image    string
	Template string
	Username string
	Password string
}

func FromEnv(skipMachine bool) (*Options, error) {
	retOptions := &Options{}

	var err error
	if !skipMachine {
		retOptions.MachineID, err = fromEnvOrError("MACHINE_ID")
		if err != nil {
			return nil, err
		}
		// prefix with devpod-
		// MachineID already includes "devpod-" prefix from DevPod
		// No need to add it again

		retOptions.MachineFolder, err = fromEnvOrError("MACHINE_FOLDER")
		if err != nil {
			return nil, err
		}
	}

	retOptions.Username, err = fromEnvOrError("UPCLOUD_USERNAME")
	if err != nil {
		return nil, err
	}
	retOptions.Password, err = fromEnvOrError("UPCLOUD_PASSWORD")
	if err != nil {
		return nil, err
	}
	retOptions.Zone, err = fromEnvOrError("UPCLOUD_ZONE")
	if err != nil {
		return nil, err
	}
	retOptions.Plan, err = fromEnvOrError("UPCLOUD_PLAN")
	if err != nil {
		return nil, err
	}
	retOptions.Storage, err = fromEnvOrError("UPCLOUD_STORAGE")
	if err != nil {
		return nil, err
	}
	retOptions.Image, err = fromEnvOrError("UPCLOUD_IMAGE")
	if err != nil {
		return nil, err
	}
	// Template is optional, so use fromEnv instead of fromEnvOrError
	retOptions.Template = os.Getenv("UPCLOUD_TEMPLATE")

	return retOptions, nil
}

// FromEnvInit is like FromEnv but only requires credentials for init command
func FromEnvInit() (*Options, error) {
	retOptions := &Options{}

	var err error
	retOptions.Username, err = fromEnvOrError("UPCLOUD_USERNAME")
	if err != nil {
		return nil, err
	}
	retOptions.Password, err = fromEnvOrError("UPCLOUD_PASSWORD")
	if err != nil {
		return nil, err
	}

	// Set defaults for other fields (not required for init)
	retOptions.Zone = os.Getenv("UPCLOUD_ZONE")
	if retOptions.Zone == "" {
		retOptions.Zone = "de-fra1"
	}
	retOptions.Plan = os.Getenv("UPCLOUD_PLAN")
	if retOptions.Plan == "" {
		retOptions.Plan = "2xCPU-4GB"
	}
	retOptions.Storage = os.Getenv("UPCLOUD_STORAGE")
	if retOptions.Storage == "" {
		retOptions.Storage = "50"
	}
	retOptions.Image = os.Getenv("UPCLOUD_IMAGE")
	if retOptions.Image == "" {
		retOptions.Image = "Ubuntu Server 22.04 LTS (Jammy Jellyfish)"
	}

	return retOptions, nil
}

func fromEnvOrError(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("couldn't find option %s in environment, please make sure %s is defined", name, name)
	}

	return val, nil
}
