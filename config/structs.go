package config

import (
	"fmt"
)

const currentMajorVersion = 1

type Gitlab struct {
	Url          string
	Token        string
	MajorVersion int
}

func (config *Gitlab) CheckValidity() error {
	// Check URL
	_, err := checkURLStr(config.Url)
	if err != nil {
		return err
	}

	// Check if token is semantically correct. The tokens validity is not checked
	if len(config.Token) < 20 {
		return fmt.Errorf("config contains token that is too short. Expected: at least 20, got %d", len(config.Token))
	}

	// Check version
	if config.MajorVersion > currentMajorVersion {
		return fmt.Errorf("config was written by a newer version of the tool")
	}

	return nil
}
