package configuration

import (
	"errors"

	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/remote"
	remotetype "github.com/glad-dev/gin/remote/type"
)

// Append adds the token to the corresponding host in the  configuration file.
// If no configuration file exists, a new one will be created.
func Append(urlStr string, token string, remoteType remotetype.Type) error {
	// Load current config
	config, err := Load()
	if err != nil && !errors.Is(ErrConfigDoesNotExist, err) {
		// Config exists, but there was some other error
		log.Error("Failed to load config", "error", err)

		return err
	}

	u, err := checkURLStr(urlStr)
	if err != nil {
		log.Error("URL is invalid.", "error", err, "url", urlStr)

		return err
	}

	// Check if a configuration with the same username and token already exists
	configLocation := -1
	for i, config := range config.Remotes {
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
		log.Error("Failed to initialize token", "error", err)

		return err
	}

	// Add new config
	if configLocation == -1 {
		// Config with given URL does not yet exist
		config.Remotes = append(config.Remotes, Remote{
			URL:     *u,
			Details: []remote.Details{rd},
		})
	} else {
		// Config with given URL exists
		config.Remotes[configLocation].Details = append(config.Remotes[configLocation].Details, rd)
	}

	// write back
	return write(config)
}
