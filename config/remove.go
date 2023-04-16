package config

import (
	"errors"
	"os"

	"gn/config/location"
	"gn/logger"
)

func Remove(wrapper *Wrapper, wrapperIndex int, detailsIndex int) error {
	// Check if index is valid
	if wrapperIndex < 0 || wrapperIndex >= len(wrapper.Remotes) {
		logger.Log.Error("Invalid wrapper index", "index", wrapperIndex, "len(remotes)", len(wrapper.Remotes))

		return errors.New("invalid wrapper index")
	}

	if detailsIndex < 0 || detailsIndex >= len(wrapper.Remotes[wrapperIndex].Details) {
		logger.Log.Error("Invalid details index", "index", detailsIndex, "len(remotes.Details)", len(wrapper.Remotes[wrapperIndex].Details))

		return errors.New("invalid details index")
	}

	if len(wrapper.Remotes[wrapperIndex].Details) == 1 {
		wrapper.Remotes = append(wrapper.Remotes[:wrapperIndex], wrapper.Remotes[wrapperIndex+1:]...)
	} else {
		wrapper.Remotes[wrapperIndex].Details = append(wrapper.Remotes[wrapperIndex].Details[:detailsIndex], wrapper.Remotes[wrapperIndex].Details[detailsIndex+1:]...)
	}

	// If there are no configs left, delete the config file
	if len(wrapper.Remotes) == 0 {
		location, err := location.Get()
		if err != nil {
			return err
		}

		return os.Remove(location)
	}

	// Write back the updated config
	return Write(wrapper)
}
