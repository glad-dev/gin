package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/BurntSushi/toml"
)

var ErrConfigDoesNotExist = errors.New("config does not exist")

func loadConfig() (*Gitlab, error) {
	fileLocation, err := getConfigLocation()
	if err != nil {
		return nil, err
	}

	// Load config
	config := Gitlab{}
	metaData, err := toml.DecodeFile(fileLocation, &config)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrConfigDoesNotExist
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
