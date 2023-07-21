package config

import (
	"errors"

	"github.com/glad-dev/gin/logger"
	"github.com/glad-dev/gin/remote"
)

// Append adds the token to the corresponding host in the  configuration file.
// If no configuration file exists, a new one will be created.
func Append(urlStr string, token string, remoteType remote.Type) error {
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
				if detail.Token == token {
					return errors.New("a configuration with the given URL and token already exists")
				}
			}
		}
	}

	rd := remote.Details{
		Type:  remoteType,
		Token: token,
	}

	err = rd.Init(u)
	if err != nil {
		logger.Log.Errorf("Failed to initialize token: %s", err)

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
