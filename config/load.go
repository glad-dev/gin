package config

import (
	"errors"
	"fmt"
	"gn/constants"
	"os"

	"github.com/BurntSushi/toml"
)

var ErrConfigDoesNotExist = errors.New("config does not exist")

// Load returns the config located at '~/.gn.toml' if it exists. If it does not exist, function returns a
// ErrConfigDoesNotExist error and an initialized General config.
func Load() (*General, error) {
	fileLocation, err := getConfigLocation()
	if err != nil {
		return nil, err
	}

	// Load config
	config := &General{}
	metaData, err := toml.DecodeFile(fileLocation, config)
	if err != nil {
		if os.IsNotExist(err) {
			return &General{
				Configs:      []GitLab{},
				MajorVersion: constants.CurrentMajorVersion,
			}, ErrConfigDoesNotExist
		}

		return nil, fmt.Errorf("could not decode config: %w", err)
	}

	// Check if the config only contains the keys we expect
	if len(metaData.Undecoded()) > 0 {
		return nil, fmt.Errorf("config contains unexpected keys: %+v", metaData.Undecoded())
	}

	err = config.CheckValidity()
	if err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	return config, nil
}
