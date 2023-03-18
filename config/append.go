package config

func Append() error {
	// Load current config
	generalConf, err := loadConfig()
	if err != nil {
		return err
	}

	// Read config from terminal
	config, err := readConfig()
	if err != nil {
		return err
	}

	// Add new config
	generalConf.Configs = append(generalConf.Configs, *config)

	// Write back
	return writeConfig(generalConf)
}
