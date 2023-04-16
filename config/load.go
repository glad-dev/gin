package config

import (
	"errors"
	"fmt"
	"os"

	"gn/config/location"
	"gn/constants"
	"gn/logger"

	"github.com/BurntSushi/toml"
)

var ErrConfigDoesNotExist = errors.New("config does not exist")

const ErrConfigDoesNotExistMsg = "No configuration exists.\nRun `" + constants.ProgramName + " config add` to add remotes"

// Load returns the config located at '~/.config/gn/gn.toml' if it exists. If it does not exist, function returns a
// ErrConfigDoesNotExist error and an initialized Wrapper config.
func Load() (*Wrapper, error) {
	fileLocation, err := location.Get()
	if err != nil {
		return nil, err
	}

	// Load config
	config := &Wrapper{}
	metaData, err := toml.DecodeFile(fileLocation, config)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Log.Infof("Found no configuration file at: %s", fileLocation)

			return &Wrapper{
				Remotes: []Remote{},
				Version: constants.ConfigVersion,
			}, ErrConfigDoesNotExist
		}

		logger.Log.Errorf("toml decode failed: %s", err)

		return nil, fmt.Errorf("could not decode config: %w", err)
	}

	// Check if the config only contains the keys we expect
	if len(metaData.Undecoded()) > 0 {
		logger.Log.Error("Config contains unexpected keys.", "invalidKeys", metaData.Undecoded())

		return nil, fmt.Errorf("config contains unexpected keys: %+v", metaData.Undecoded())
	}

	err = config.CheckValidity()
	if err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	err = config.Colors.setColors()
	if err != nil {
		return nil, err
	}

	return config, nil
}
