package configuration

import (
	"errors"

	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/remote"
	rt "github.com/glad-dev/gin/remote/type"
)

// ErrUpdateSameValues is returned if Update was called with the same url and token that is already stored in the
// configuration file.
var ErrUpdateSameValues = errors.New("called update config with existing values")

// Update updates the token/url and checks the token's validity.
func Update(config *Config, configIndex int, detailsIndex int, url string, token string) error {
	if configIndex < 0 || configIndex >= len(config.Remotes) {
		log.Error("Config index is invalid.", "index", configIndex, "len(remotes)", len(config.Remotes))

		return errors.New("update: invalid config index")
	}

	if detailsIndex < 0 || detailsIndex >= len(config.Remotes[configIndex].Details) {
		log.Error("Details index is invalid.", "index", detailsIndex, "len(details)", len(config.Remotes[configIndex].Details))

		return errors.New("update: invalid details index")
	}

	u, err := parseURLStr(url)
	if err != nil {
		return err
	}

	// Check if there are any changes
	old := config.Remotes[configIndex]
	if old.URL == *u {
		for _, detail := range old.Details {
			if detail.Token == token {
				log.Warn("Attempted to update the remote '%s' with the same token", u.String())

				return ErrUpdateSameValues
			}
		}
	}

	rd := remote.Details{
		Token: token,
	}

	var remoteType rt.Type
	if u.Host == "github.com" {
		remoteType = rt.Github
	} else {
		remoteType = rt.Gitlab
	}

	err = rd.Init(u, remoteType)
	if err != nil {
		log.Error("Failed to initialize token", "error", err)

		return err
	}

	config.Remotes[configIndex].Type = remoteType
	config.Remotes[configIndex].Details[detailsIndex] = rd

	return write(config)
}
