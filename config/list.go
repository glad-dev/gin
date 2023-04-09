package config

import (
	"fmt"

	"gn/tui/style"
)

func List() error {
	// Load config
	wrapper, err := Load()
	if err != nil {
		return err
	}

	configLocation, err := getConfigLocation()
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Print(style.Title.Render(fmt.Sprintf("The configuration file at '%s' contains the following remotes:", configLocation)))
	for i, config := range wrapper.Configs {
		fmt.Print(style.List.Render(fmt.Sprintf("%d) %s", i+1, config.URL.String())))

		for k, detail := range config.Details {
			fmt.Print(style.ListDetails.Render(fmt.Sprintf(
				"%d.%d) Username: '%s' - Token name: '%s'", // TODO: Improve output
				i+1, k+1,
				detail.Username,
				detail.TokenName,
			)))
		}
	}
	fmt.Print("\n\n")

	return nil
}
