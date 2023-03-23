package config

import (
	"errors"
)

func AppendOnce(urlStr string, token string) error {
	if firstCall {
		firstCall = false
		errPrevious = Append(urlStr, token)
	}

	return errPrevious
}

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

	config := GitLab{
		URL:   *u,
		Token: token,
	}

	// Add new config
	generalConf.Configs = append(generalConf.Configs, config)

	// Write back
	return writeConfig(generalConf)
}
