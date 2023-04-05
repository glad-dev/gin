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
	for _, config := range wrapper.Configs {
		if config.URL == *u && config.Token == token {
			return fmt.Errorf("a configuration with the given URL and token already exists")
		}
	}

	lab := GitLab{
		URL:   *u,
		Token: token,
	}

	err = lab.Init()
	if err != nil {
		return err
	}

	// Add new config
	wrapper.Configs = append(wrapper.Configs, lab)

	// Write back
	return writeConfig(wrapper)
}
