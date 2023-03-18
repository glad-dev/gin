package config

import "fmt"

func List() error {
	// Load config
	generalConfig, err := loadConfig()
	if err != nil {
		return err
	}

	configLocation, err := getConfigLocation()
	if err != nil {
		return err
	}

	fmt.Printf("The configuration file at '%s' contains data for the following URLs:\n", configLocation)
	for i, config := range generalConfig.Configs {
		fmt.Printf("%d) %s\n", i+1, config.URL.String())
	}

	return nil
}
