package configuration

import (
	"fmt"

	"github.com/glad-dev/gin/configuration/location"
	"github.com/glad-dev/gin/style"
)

// List prints all stored remotes.
func List() error {
	// Load config
	config, err := Load()
	if err != nil {
		return err
	}

	configLocation, err := location.Get()
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(style.Title.Render(fmt.Sprintf("The configuration file at '%s' contains the following remotes:", configLocation)))
	for i, conf := range config.Remotes {
		fmt.Println(style.PrintOnlyList.Render(fmt.Sprintf("%d) %s", i+1, conf.URL.String())))

		for k, detail := range conf.Details {
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
