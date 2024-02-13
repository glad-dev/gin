package config

import (
	"errors"
	"os"

	"github.com/glad-dev/gin/config/location"
	"github.com/glad-dev/gin/log"
)

// Remove removes the token/url combination at the passed indices.
func Remove(wrapper *Wrapper, wrapperIndex int, detailsIndex int) error {
	// Check if index is valid
	if wrapperIndex < 0 || wrapperIndex >= len(wrapper.Remotes) {
		log.Error("Invalid wrapper index.", "index", wrapperIndex, "len(remotes)", len(wrapper.Remotes))

		return errors.New("invalid wrapper index")
	}

	if detailsIndex < 0 || detailsIndex >= len(wrapper.Remotes[wrapperIndex].Details) {
		log.Error("Invalid details index.", "index", detailsIndex, "len(remotes.Details)", len(wrapper.Remotes[wrapperIndex].Details))

		return errors.New("invalid details index")
	}

	if len(wrapper.Remotes[wrapperIndex].Details) == 1 {
		wrapper.Remotes = append(wrapper.Remotes[:wrapperIndex], wrapper.Remotes[wrapperIndex+1:]...)
	} else {
		wrapper.Remotes[wrapperIndex].Details = append(wrapper.Remotes[wrapperIndex].Details[:detailsIndex], wrapper.Remotes[wrapperIndex].Details[detailsIndex+1:]...)
	}

	// If there are no configs left, delete the config file
	if len(wrapper.Remotes) == 0 {
		loc, err := location.Get()
		if err != nil {
			return err
		}

		return os.Remove(loc)
	}

	// write back the updated config
	return write(wrapper)
}
