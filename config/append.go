package config

import (
	"errors"

	"gn/config/remote"
	"gn/logger"
)

func Append(urlStr string, token string) error {
	// Load current config
	wrapper, err := Load()
	if err != nil && !errors.Is(ErrConfigDoesNotExist, err) {
		// Config exists, but there was some other error
		logger.Log.Errorf("Failed to load config: %s", err)

		return err
	}

	u, err := checkURLStr(urlStr)
	if err != nil {
		logger.Log.Error("URL is invalid.", "error", err, "url", urlStr)

		return err
	}

	// Check if a configuration with the same username and token already exists
	configLocation := -1
	for i, config := range wrapper.Remotes {
		if config.URL == *u {
			configLocation = i

			for _, detail := range config.Details {
				if detail.GetToken() == token {
					return errors.New("a configuration with the given URL and token already exists")
				}
			}
		}
	}

	var rd remote.Details
	if u.Host == "github.com" {
		rd = remote.GitHubDetails{
			Token: token,
		}
	} else {
		rd = remote.GitLabDetails{
			Token: token,
		}
	}

	rd, err = rd.Init(u)
	if err != nil {
		logger.Log.Errorf("Failed to initalize token: %s", err)

		return err
	}

	// Add new config
	if configLocation == -1 {
		// Config with given URL does not yet exist
		wrapper.Remotes = append(wrapper.Remotes, Remote{
			URL:     *u,
			Details: []remote.Details{rd},
		})
	} else {
		// Config with given URL exists
		wrapper.Remotes[configLocation].Details = append(wrapper.Remotes[configLocation].Details, rd)
	}

	// Write back
	return Write(wrapper)
}
