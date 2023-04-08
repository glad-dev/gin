package config

import (
	"errors"
	"fmt"
)

func Append(urlStr string, token string) error {
	// Load current config
	wrapper, err := Load()
	if err != nil && !errors.Is(ErrConfigDoesNotExist, err) {
		// Config exists, but there was some other error
		return err
	}

	u, err := checkURLStr(urlStr)
	if err != nil {
		return err
	}

	// Check if a configuration with the same username and token already exists
	configLocation := -1
	for i, config := range wrapper.Configs {
		if config.URL == *u {
			configLocation = i

			for _, detail := range config.Details {
				if detail.Token == token {
					return fmt.Errorf("a configuration with the given URL and token already exists")
				}
			}
		}
	}

	rd := RepoDetails{
		Token: token,
	}

	err = rd.Init(u)
	if err != nil {
		return err
	}

	// Add new config
	if configLocation == -1 {
		// Config with given URL does not yet exist
		wrapper.Configs = append(wrapper.Configs, Repo{
			URL:     *u,
			Details: []RepoDetails{rd},
		})
	} else {
		// Config with given URL exists
		wrapper.Configs[configLocation].Details = append(wrapper.Configs[configLocation].Details, rd)
	}

	// Write back
	return writeConfig(wrapper)
}
