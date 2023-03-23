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
	fmt.Print(style.TitleStyle.Render(fmt.Sprintf("The configuration file at '%s' contains data for the following URLs:", configLocation)))
	for i, config := range generalConfig.Configs {
		fmt.Print(style.ListStyle.Render(fmt.Sprintf("%d) %s", i+1, config.URL.String())))
	}
	fmt.Print("\n\n")

	return nil
}
