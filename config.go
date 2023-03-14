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

func readConfig() (*GitlabConfig, error) {
	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("could not read gitlab config: %w", err)
	}

	// Load config
	config := GitlabConfig{}
	_, err = toml.DecodeFile(path.Join(usr.HomeDir, ".gn.toml"), &config)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("the config file at \"~/.gn.toml\" does not exist")
		}

		return nil, fmt.Errorf("could not decode config: %w", err)
	}

	return &config, nil
}
