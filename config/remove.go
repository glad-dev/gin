package config

import (
	"errors"
	"os"
)

func Remove(wrapper *Wrapper, index int) error {
	// Check if index is valid
	if index < 0 || index >= len(wrapper.Configs) {
		return errors.New("invalid index")
	}

	// Remove the selected config
	wrapper.Configs = append(wrapper.Configs[:index], wrapper.Configs[index+1:]...)

	// If there are no configs, delete the config file
	if len(wrapper.Configs) == 0 {
		location, err := getConfigLocation()
		if err != nil {
			return err
		}

		return os.Remove(location)
	}

	// Write back the updated config
	return writeConfig(wrapper)
}
