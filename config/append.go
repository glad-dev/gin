package config

func Append() error {
	// Load current config
	generalConf, err := Load()
	if err != nil {
		return err
	}

	// Read config from terminal
	config, err := readConfigFromStdIn()
	if err != nil {
		return err
	}

	// Add new config
	generalConf.Configs = append(generalConf.Configs, *config)

	// Write back
	return writeConfig(generalConf)
}
