package config

func Remove() error {
	// Load current config
	generalConfig, err := loadConfig()
	if err != nil {
		return err
	}

	index, err := selectExistingConfigs(generalConfig.Configs)
	if err != nil {
		return err
	}

	// Remove the selected config
	generalConfig.Configs = append(generalConfig.Configs[:index], generalConfig.Configs[index+1:]...)

	// Write back the updated config
	return writeConfig(generalConfig)
}
