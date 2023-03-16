package config

import (
	"errors"
	"os"
)

func Get() (*Gitlab, error) {
	config, err := loadConfig()
	if err != nil {
		if !errors.Is(err, ErrConfigDoesNotExist) {
			// Error is NOT about the config not existing
			return nil, err
		}

		config, err = readConfigFromStdin()
		if err != nil {
			if errors.Is(err, ErrDontCreateConfig) {
				os.Exit(0)
			}

			return nil, err
		}

		err = writeConfig(config)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}
