package config

func AppendOnce(urlStr string, token string) error {
	if firstCall {
		firstCall = false
		errPrevious = Append(urlStr, token)
	}

	return errPrevious
}

func Append(urlStr string, token string) error {
	// Load current config
	generalConf, err := Load() // ToDo: Deal with adding entry to empty config
	if err != nil {
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
