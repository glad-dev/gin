package config

import (
	"errors"
	"fmt"
	"net/url"

	"gn/constants"
)

type General struct {
	Configs      []GitLab
	MajorVersion int
}

type GitLab struct {
	URL   url.URL
	Token string
}

var ErrNoMatchingConfig = errors.New("no matching config was found")

func (config *General) CheckValidity() error {
	if len(config.Configs) == 0 {
		return errors.New("config file does not contain []GitLab")
	}

	// Check version
	if config.MajorVersion > constants.CurrentMajorVersion {
		return fmt.Errorf("config was written by a newer version of the tool")
	}

	for _, singleConfig := range config.Configs {
		// Check URL
		_, err := checkURLStr(singleConfig.URL.String())
		if err != nil {
			return err
		}

		// Check if token is semantically correct. The tokens validity is not checked
		if len(singleConfig.Token) < 20 {
			return fmt.Errorf("config contains token that is too short. Expected: at least 20, got %d", len(singleConfig.Token))
		}
	}

	return nil
}

func (config *General) GetMatchingConfig(rawURL string) (*GitLab, error) {
	// Convert url string to url.Url
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}

	for _, lab := range config.Configs {
		if u.Host == lab.URL.Host {
			return &lab, err
		}
	}

	return nil, ErrNoMatchingConfig
}
