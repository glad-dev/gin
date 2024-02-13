package configuration

import (
	"errors"
	"os"

	"github.com/glad-dev/gin/configuration/location"
	"github.com/glad-dev/gin/log"
)

// Remove removes the token/url combination at the passed indices.
func Remove(config *Config, configIndex int, detailsIndex int) error {
	// Check if index is valid
	if configIndex < 0 || configIndex >= len(config.Remotes) {
		log.Error("Invalid config index.", "index", configIndex, "len(remotes)", len(config.Remotes))

		return errors.New("invalid config index")
	}

	if detailsIndex < 0 || detailsIndex >= len(config.Remotes[configIndex].Details) {
		log.Error("Invalid details index.", "index", detailsIndex, "len(remotes[index].Details)", len(config.Remotes[configIndex].Details))

		return errors.New("invalid details index")
	}

	if len(config.Remotes[configIndex].Details) == 1 {
		config.Remotes = append(config.Remotes[:configIndex], config.Remotes[configIndex+1:]...)
	} else {
		config.Remotes[configIndex].Details = append(config.Remotes[configIndex].Details[:detailsIndex], config.Remotes[configIndex].Details[detailsIndex+1:]...)
	}

	// If there are no configs left, delete the config file
	if len(config.Remotes) == 0 {
		loc, err := location.Get()
		if err != nil {
			return err
		}

		return os.Remove(loc)
	}

	// write back the updated config
	return write(config)
}
