package configuration

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/glad-dev/gin/configuration/location"
	"github.com/glad-dev/gin/constants"
	"github.com/glad-dev/gin/log"

	"github.com/BurntSushi/toml"
)

// write checks if the passed config is valid and writes it to ~/.config/gn/gn.toml.
func write(config *Config) error {
	openConfig := func(fileLocation string) (*os.File, error) {
		return os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	}

	config.Version = constants.ConfigVersion

	err := config.CheckValidity()
	if err != nil {
		return fmt.Errorf("passed config is invalid: %w", err)
	}

	buf := new(bytes.Buffer)
	err = toml.NewEncoder(buf).Encode(config)
	if err != nil {
		log.Error("Failed to encode config", "error", err)

		return fmt.Errorf("could not encode config: %w", err)
	}

	fileLocation, err := location.Get()
	if err != nil {
		return err
	}

	f, err := openConfig(fileLocation)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Error("Failed to open config file", "error", err)

			return fmt.Errorf("could not open config file: %w", err)
		}

		// Create config directory
		err = location.CreateDir()
		if err != nil {
			log.Error("Failed to create config directory", "error", err)

			return fmt.Errorf("could not create config directory: %w", err)
		}

		// Attempt to create the config file
		f, err = openConfig(fileLocation)
		if err != nil {
			log.Error("Failed to open newly created config file", "error", err)

			return fmt.Errorf("could not open config file: %w", err)
		}
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		log.Error("Failed to write config to file", "error", err)

		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}
