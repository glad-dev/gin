package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func writeConfig(config *Gitlab) error {
	config.MajorVersion = currentMajorVersion

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

	f, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("could not open config file: %w", err)
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}
