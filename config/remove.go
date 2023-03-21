package config

import (
	"errors"
	"os"
)

var removeOnceAlreadyCalled = false

func RemoveOnce(generalConfig *General, index int) error {
	if removeOnceAlreadyCalled {
		return nil
	}
	removeOnceAlreadyCalled = true

	return Remove(generalConfig, index)
}

func Remove(generalConfig *General, index int) error {
	// Check if index is valid
	if index < 0 || index >= len(generalConfig.Configs) {
		return errors.New("invalid index")
	}

	// Remove the selected config
	generalConfig.Configs = append(generalConfig.Configs[:index], generalConfig.Configs[index+1:]...)

	// If there are no GitLab configs, delete the config file
	if len(generalConfig.Configs) == 0 {
		location, err := getConfigLocation()
		if err != nil {
			return err
		}

		return os.Remove(location)
	}

	// Write back the updated config
	return writeConfig(generalConfig)
}
