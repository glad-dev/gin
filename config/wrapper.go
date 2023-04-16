package config

import (
	"errors"
	"fmt"

	"gn/config/remote"
	"gn/constants"
	"gn/logger"
	"gn/repo"
)

type Wrapper struct {
	Colors  Colors
	Remotes []Remote
	Version uint8
}

var ErrNoMatchingConfig = errors.New("no matching config was found")

func (config *Wrapper) CheckValidity() error {
	if len(config.Remotes) == 0 {
		logger.Log.Error("Config does not contain any remotes.")

		return errors.New("config file does not contain remotes")
	}

	// Check version
	if config.Version > constants.ConfigVersion {
		logger.Log.Error("Config has newer version than the program.", "configVersion", config.Version, "expectedVersion", constants.ConfigVersion)

		return fmt.Errorf("config was written by a newer version of the tool")
	}

	// Check configs
	for _, remote := range config.Remotes {
		err := remote.checkSemantics()
		if err != nil {
			logger.Log.Error("Invalid remote.", "error", err, "remote", remote)

			return err
		}
	}

	// Check colors
	err := config.Colors.CheckValidity()
	if err != nil {
		logger.Log.Error("")

		return err
	}

	return nil
}

func (config *Wrapper) GetMatchingConfig(details []repo.Details) (*remote.Match, string, error) {
	if len(details) == 0 {
		logger.Log.Error("No details passed.")

		return nil, "", errors.New("no details passed")
	}

	for _, detail := range details {
		for _, conf := range config.Remotes {
			if conf.URL.Host == detail.URL.Host {
				match, err := conf.ToMatch()
				if err != nil {
					return nil, "", err
				}

				return match, detail.ProjectPath, nil
			}
		}
	}

	return nil, "", ErrNoMatchingConfig
}
