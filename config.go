package main

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path"
	"reflect"

	"github.com/BurntSushi/toml"
)

type GitlabConfig struct {
	Url          string
	Token        string
	MajorVersion int
}

const currentMajorVersion = 1

func (config *GitlabConfig) CheckValidity() error {
	// Check URL
	u, err := url.ParseRequestURI(config.Url)
	if err != nil {
		return fmt.Errorf("config contains invalid URL: %w", err)
	}

	if !u.IsAbs() {
		return fmt.Errorf("config contains URL that is not absolute: %s", config.Url)
	}

	// Check if token is semantically correct. The tokens validity is not checked
	if len(config.Token) < 20 {
		return fmt.Errorf("config contains token that is too short. Expected: at least 20, got %d", len(config.Token))
	}

	if config.MajorVersion > currentMajorVersion {
		return fmt.Errorf("config was written by a newer version of the tool")
	}

	return nil
}

func writeConfig(config *GitlabConfig) error {
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

func readConfig() (*GitlabConfig, error) {
	fileLocation, err := getConfigLocation()
	if err != nil {
		return nil, err
	}

	// Load config
	config := GitlabConfig{}
	metaData, err := toml.DecodeFile(fileLocation, &config)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file at '~/.gn.toml' does not exist")
		}

		return nil, fmt.Errorf("could not decode config: %w", err)
	}

	// Check if the config only contains the keys we expect
	if len(metaData.Undecoded()) > 0 {
		return nil, fmt.Errorf("config contains unexpected keys: %+v", metaData.Undecoded())
	}

	// Check if the config contains all the keys we need
	expected := reflect.ValueOf(config).NumField()
	if len(metaData.Keys()) != expected {
		return nil, fmt.Errorf("config contains an invalid amount of keys. Expect %d, got %d", expected, len(metaData.Keys()))
	}

	err = config.CheckValidity()
	if err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	return &config, nil
}

func getConfigLocation() (string, error) {
	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("could not get current user: %w", err)
	}

	return path.Join(usr.HomeDir, ".gn.toml"), nil
}
