package config

import (
	"errors"
	"fmt"

	"github.com/glad-dev/gin/constants"
	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/remote/match"
	"github.com/glad-dev/gin/repo"
)

// Wrapper contains the Colors configurations, a Version number and a list of Remotes.
type Wrapper struct {
	Colors  Colors
	Remotes []Remote
	Version uint8
}

// CheckValidity checks if the wrapper is valid. It checks the wrapper's version number, colors, and remotes.
func (config *Wrapper) CheckValidity() error {
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

// GetMatchingConfig searches the wrapper's Remotes and returns a remote.Match if a Remote has the same  URL as one of
// the passed repo.Details.
func (config *Wrapper) GetMatchingConfig(details []repo.Details) (*match.Match, string, error) {
	if len(details) == 0 {
		log.Error("No details passed.")

		return nil, "", errors.New("no details passed")
	}

	for _, detail := range details {
		for _, conf := range config.Remotes {
			if conf.URL.Host == detail.URL.Host {
				m, err := conf.ToMatch()
				if err != nil {
					return nil, "", err
				}

				return m, detail.ProjectPath, nil
			}
		}
	}

	// No match => Mock up a config
	return &match.Match{
		URL:      details[0].URL,
		Token:    "",
		Username: "",
	}, details[0].ProjectPath, nil
}
