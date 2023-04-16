package config

import (
	"fmt"

	"gn/config/location"
	"gn/style"
)

func List() error {
	// Load config
	wrapper, err := Load()
	if err != nil {
		return err
	}

	configLocation, err := location.Get()
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(style.Title.Render(fmt.Sprintf("The configuration file at '%s' contains the following remotes:", configLocation)))
	for i, config := range wrapper.Configs {
		fmt.Println(style.List.Render(fmt.Sprintf("%d) %s", i+1, config.URL.String())))

		for k, detail := range config.Details {
			fmt.Println(style.ListDetails.Render(fmt.Sprintf(
				"%d.%d) Username: '%s' - Token name: '%s'", // TODO: Improve output
				i+1, k+1,
				detail.Username,
				detail.TokenName,
			)))
		}
	}
	fmt.Print("\n")

	return nil
}
