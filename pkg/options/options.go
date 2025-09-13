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
		retOptions.MachineID = "devpod-" + retOptions.MachineID

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

	return retOptions, nil
}

func fromEnvOrError(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("couldn't find option %s in environment, please make sure %s is defined", name, name)
	}

	return val, nil
}
