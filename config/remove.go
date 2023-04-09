package config

import (
	"errors"
	"os"
)

func Remove(wrapper *Wrapper, wrapperIndex int, detailsIndex int) error {
	// Check if index is valid
	if wrapperIndex < 0 || wrapperIndex >= len(wrapper.Configs) {
		return errors.New("invalid wrapper index")
	}

	if detailsIndex < 0 || detailsIndex >= len(wrapper.Configs[wrapperIndex].Details) {
		return errors.New("invalid details index")
	}

	if len(wrapper.Configs[wrapperIndex].Details) == 1 {
		wrapper.Configs = append(wrapper.Configs[:wrapperIndex], wrapper.Configs[wrapperIndex+1:]...)
	} else {
		wrapper.Configs[wrapperIndex].Details = append(wrapper.Configs[wrapperIndex].Details[:detailsIndex], wrapper.Configs[wrapperIndex].Details[detailsIndex+1:]...)
	}

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
