package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/glad-dev/gin/config/location"
	"github.com/glad-dev/gin/constants"
	"github.com/glad-dev/gin/logger"

	"github.com/BurntSushi/toml"
)

// Write checks if the passed config is valid and writes it to ~/.config/gn/gn.toml.
func Write(config *Wrapper) error {
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
		logger.Log.Errorf("Failed to encode config: %s", err)

		return fmt.Errorf("could not encode config: %w", err)
	}

	fileLocation, err := location.Get()
	if err != nil {
		return err
	}

	f, err := openConfig(fileLocation)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Log.Errorf("Failed to open config file: %s", err)

			return fmt.Errorf("could not open config file: %w", err)
		}

		// Create config directory
		err = location.CreateDir()
		if err != nil {
			logger.Log.Errorf("Failed to create config directory: %s", err)

			return fmt.Errorf("could not create config directory: %w", err)
		}

		// Attempt to create the config file
		f, err = openConfig(fileLocation)
		if err != nil {
			logger.Log.Errorf("Failed to open newly created config file: %s", err)

			return fmt.Errorf("could not open config file: %w", err)
		}
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		logger.Log.Errorf("Failed to write config to file: %s", err)

		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}
