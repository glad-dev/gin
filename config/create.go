package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"gn/constants"
)

var ErrDontCreateConfig = errors.New("user does not want to create the config")

func readConfigFromStdin() (*Wrapper, error) {
	path, err := getConfigLocation()
	if err != nil {
		return nil, err
	}

	fmt.Printf("No config file was found at: %s\n", path)
	fmt.Printf("Do you want to create it (y/n)? ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	input := strings.ToLower(scanner.Text())
	switch input {
	case "n":
		return nil, ErrDontCreateConfig
	case "yes":
		fallthrough
	case "y":
		break
	default:
		return nil, fmt.Errorf("invalid input. Expected 'y' or 'n', got '%s'", input)
	}

	config, err := readConfigFromStdIn()
	if err != nil {
		return nil, err
	}

	return &Wrapper{
		MajorVersion: constants.CurrentMajorVersion,
		Configs: []GitLab{
			*config,
		},
	}, nil
}
