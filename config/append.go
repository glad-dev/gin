package config

import (
	"errors"
)

func Append(urlStr string, token string) error {
	// Load current config
	generalConf, err := Load()
	if err != nil && !errors.Is(ErrConfigDoesNotExist, err) {
		// Config exists, but there was some other error
		return err
	}

	u, err := checkURLStr(urlStr)
	if err != nil {
		return err
	}

	lab := GitLab{
		URL:   *u,
		Token: token,
	}

	err = lab.Init()
	if err != nil {
		return err
	}

	// Add new config
	generalConf.Configs = append(generalConf.Configs, lab)

	// Write back
	return writeConfig(generalConf)
}
