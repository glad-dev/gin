package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"gn/constants"

	"github.com/BurntSushi/toml"
)

func Write(config *Wrapper) error {
	config.ConfigVersion = constants.ConfigVersion

	err := config.CheckValidity()
	if err != nil {
		return fmt.Errorf("passed config is invalid: %w", err)
	}

	buf := new(bytes.Buffer)
	err = toml.NewEncoder(buf).Encode(config)
	if err != nil {
		return fmt.Errorf("could not encode config: %w", err)
	}

	fileLocation, err := getConfigLocation()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("could not open config file: %w", err)
		}

		// Create config directory
		err = createConfigDir()
		if err != nil {
			return fmt.Errorf("could not create config directory: %w", err)
		}

		// Attempt to create the config file
		f, err = os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		if err != nil {
			return fmt.Errorf("could not open config file: %w", err)
		}
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}
