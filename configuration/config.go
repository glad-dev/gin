package configuration

import (
	"errors"
	"fmt"

	"github.com/glad-dev/gin/constants"
	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/remote/match"
	rt "github.com/glad-dev/gin/remote/type"
	"github.com/glad-dev/gin/repository"
)

// Config contains the Colors configurations, a Version number and a list of Remotes.
type Config struct {
	Colors  Colors
	Remotes []Remote
	Version uint8
}

// CheckValidity checks if the config is valid. It checks the config's version number, colors, and remotes.
func (config *Config) CheckValidity() error {
	if len(config.Remotes) == 0 {
		log.Error("Config does not contain any remotes.")

		return errors.New("config file does not contain remotes")
	}

	// Check version
	if config.Version > constants.ConfigVersion {
		log.Error("Config has newer version than the program.", "configVersion", config.Version, "expectedVersion", constants.ConfigVersion)

		return fmt.Errorf("config was written by a newer version of the tool")
	}

	// Check configs
	for _, r := range config.Remotes {
		err := r.checkSemantics()
		if err != nil {
			log.Error("Invalid remote.", "error", err, "remote", r)

			return err
		}
	}

	// Check colors
	err := config.Colors.CheckValidity()
	if err != nil {
		log.Error("")

		return err
	}

	return nil
}

// GetMatchingConfig searches the config's Remotes and returns a remote.Match if a Remote has the same  URL as one of
// the passed repository.Details.
func (config *Config) GetMatchingConfig(details []repository.Details) (*match.Match, string, error) {
	if len(details) == 0 {
		log.Error("No details passed.")

		return nil, "", errors.New("no details passed")
	}

	for _, detail := range details {
		for _, remote := range config.Remotes {
			if remote.URL.Host == detail.URL.Host {
				m, err := remote.ToMatch()
				if err != nil {
					return nil, "", err
				}

				return m, detail.ProjectPath, nil
			}
		}
	}

	// No match => Mock up a config
	t := rt.Gitlab
	if details[0].URL.Host == "github.com" {
		t = rt.Github
	}

	return &match.Match{
		URL:      details[0].URL,
		Token:    "",
		Type:     t,
		Username: "",
	}, details[0].ProjectPath, nil
}
