package config

import (
	"fmt"

	style "gn/tui/style/config"
)

func List() error {
	// Load config
	generalConfig, err := Load()
	if err != nil {
		return err
	}

	configLocation, err := getConfigLocation()
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Print(style.Title.Render(fmt.Sprintf("The configuration file at '%s' contains the following remotes:", configLocation)))
	for i, config := range generalConfig.Configs {
		fmt.Print(style.List.Render(fmt.Sprintf("%d) %s", i+1, config.URL.String())))
	}
	fmt.Print("\n\n")

	return nil
}
