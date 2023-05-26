package config

import (
	"fmt"

	"github.com/glad-dev/gin/config/location"
	"github.com/glad-dev/gin/style"
)

// List prints all stored tokens.
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
	for i, config := range wrapper.Remotes {
		fmt.Println(style.PrintOnlyList.Render(fmt.Sprintf("%d) %s", i+1, config.URL.String())))

		for k, detail := range config.Details {
			fmt.Println(style.ListDetails.Render(fmt.Sprintf(
				"%d.%d) Username: '%s' - Token name: '%s'", // TODO: Improve output
				i+1, k+1,
				detail.GetUsername(),
				detail.GetTokenName(),
			)))
		}
	}
	fmt.Print("\n")

	return nil
}
