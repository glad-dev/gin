package config

import (
	"errors"
	"fmt"

	"gn/constants"
	"gn/repo"
)

type Wrapper struct {
	Colors        Colors
	Configs       []Remote
	ConfigVersion int
}

var ErrNoMatchingConfig = errors.New("no matching config was found")

func (config *Wrapper) CheckValidity() error {
	if len(config.Configs) == 0 {
		return errors.New("config file does not contain []Match")
	}

	// Check version
	if config.ConfigVersion > constants.ConfigVersion {
		return fmt.Errorf("config was written by a newer version of the tool")
	}

	// Check configs
	for _, singleConfig := range config.Configs {
		err := singleConfig.CheckSemantics()
		if err != nil {
			return err
		}
	}

	// Check colors
	err := config.Colors.CheckValidity()
	if err != nil {
		return err
	}

	return nil
}

func (config *Wrapper) GetMatchingConfig(details []repo.Details) (*Match, string, error) {
	if len(details) == 0 {
		return nil, "", errors.New("no details passed")
	}

	for _, detail := range details {
		for _, conf := range config.Configs {
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
