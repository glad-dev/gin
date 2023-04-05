package config

import (
	"errors"
	"fmt"

	"gn/constants"
	"gn/repo"
)

type Wrapper struct {
	Configs       []GitLab
	ConfigVersion int
}

var ErrNoMatchingConfig = errors.New("no matching config was found")

func (config *Wrapper) CheckValidity() error {
	if len(config.Configs) == 0 {
		return errors.New("config file does not contain []GitLab")
	}

	// Check version
	if config.ConfigVersion > constants.ConfigVersion {
		return fmt.Errorf("config was written by a newer version of the tool")
	}

	for _, singleConfig := range config.Configs {
		// Check URL
		_, err := checkURLStr(singleConfig.URL.String())
		if err != nil {
			return err
		}

		// Check if token is semantically correct. The tokens validity is not checked
		if len(singleConfig.Token) < 20 { // TODO: Get actual sizes
			return fmt.Errorf("config contains token that is too short. Expected: at least 20, got %d", len(singleConfig.Token))
		}

		if len(singleConfig.Username) == 0 {
			return fmt.Errorf("config contains empty username")
		}
	}

	return nil
}

func (config *Wrapper) GetMatchingConfig(details []repo.Details) (*GitLab, string, error) {
	if len(details) == 0 {
		return nil, "", errors.New("no details passed")
	}

	for _, detail := range details {
		for _, lab := range config.Configs {
			if lab.URL.Host == detail.URL.Host {
				return &lab, detail.ProjectPath, nil
			}
		}
	}

	return nil, "", ErrNoMatchingConfig
}
