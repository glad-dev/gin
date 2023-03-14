package main

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/BurntSushi/toml"
)

type GitlabConfig struct {
	Url   string
	Token string
}

func getConfigLocation() (string, error) {
	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("could not get current user: %w", err)
	}

	return path.Join(usr.HomeDir, ".gn.toml"), nil
}

func readConfig() (*GitlabConfig, error) {
	fileLocation, err := getConfigLocation()
	if err != nil {
		return nil, err
	}

	// Load config
	config := GitlabConfig{}
	metaData, err := toml.DecodeFile(fileLocation, &config)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file at \"~/.gn.toml\" does not exist")
		}

		return nil, fmt.Errorf("could not decode config: %w", err)
	}

	// Check if the config only contains the keys we expect
	if len(metaData.Undecoded()) > 0 {
		return nil, fmt.Errorf("config contains unexpected keys: %+v", metaData.Undecoded())
	}

	return &config, nil
}
