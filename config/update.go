package config

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func UpdateToken() error {
	// Load the current config
	wrapper, err := Load()
	if err != nil {
		return err
	}

	index, err := selectExistingConfigs(wrapper.Configs)
	if err != nil {
		return err
	}

	// Read token
	fmt.Printf("Enter the API token (input is hidden): ")
	token, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Printf("\n")
	if err != nil {
		return err
	}

	wrapper.Configs[index].Token = string(token)

	// Write back updated config
	return writeConfig(wrapper)
}
