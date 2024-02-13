package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/glad-dev/gin/config/location"
	"github.com/glad-dev/gin/constants"
	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/style"

	"github.com/BurntSushi/toml"
)

// ErrConfigDoesNotExist is returned by Load if no configuration file exists.
var ErrConfigDoesNotExist = errors.New("config does not exist")

const ErrConfigDoesNotExistMsg = "No configuration exists.\nRun `gin config add` to add remotes"

// Load returns the config located at '~/.config/gin/config.toml' if it exists. If it does not exist, function returns a
// ErrConfigDoesNotExist error and an initialized Wrapper config.
func Load() (*Wrapper, error) {
	fileLocation, err := location.Get()
	if err != nil {
		return nil, err
	}

	// Load config
	wrap := &Wrapper{}
	metaData, err := toml.DecodeFile(fileLocation, wrap)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info("Found no configuration file", "location", fileLocation)

			return &Wrapper{
				Remotes: []Remote{},
				Version: constants.ConfigVersion,
			}, ErrConfigDoesNotExist
		}

		log.Error("toml decode failed", "error", err)

		return nil, fmt.Errorf("could not decode config: %w", err)
	}

	// Check if the config only contains the keys we expect
	if len(metaData.Undecoded()) > 0 {
		log.Error("Config contains unexpected keys.", "invalidKeys", metaData.Undecoded())

		return nil, fmt.Errorf("config contains unexpected keys: %+v", metaData.Undecoded())
	}

	err = wrap.CheckValidity()
	if err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	err = wrap.Colors.setColors()
	if err != nil {
		return nil, err
	}

	style.UpdateColors()

	return wrap, nil
}
