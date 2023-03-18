package config

import (
	"errors"
	"fmt"
	"gn/constants"
	"net/url"
)

type General struct {
	MajorVersion int
	Configs      []GitLab
}

type GitLab struct {
	Url   url.URL
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
		_, err := checkURLStr(singleConfig.Url.String())
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

func (config *General) GetMatchingConfig(rawUrl string) (*GitLab, error) {
	// Convert url string to url.Url
	u, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return nil, err
	}

	for _, lab := range config.Configs {
		if u.Host == lab.Url.Host {
			return &lab, err
		}
	}

	return nil, ErrNoMatchingConfig
}
