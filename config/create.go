package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

var ErrDontCreateConfig = errors.New("user does not want to create the config")

func readConfigFromStdin() (*Gitlab, error) {
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

	config := Gitlab{
		Url:          "",
		Token:        "",
		MajorVersion: currentMajorVersion,
	}

	fmt.Printf("What is the URL? ")
	scanner.Scan()
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	config.Url = scanner.Text()
	u, err := checkURLStr(config.Url)
	if err != nil {
		return nil, err
	}

	fullURL := u.JoinPath("-/profile/personal_access_tokens")
	rest := "?name=git-navigator&scopes=api,read_api,read_user" // Can't be added with url.JoinPath since that escapes the '?'
	fmt.Printf("Go to %s%s to create an API key with the permissions api, read_api and read_user\n", fullURL.String(), rest)

	fmt.Printf("Enter the API token (input is hidden): ")
	token, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Printf("\n")
	if err != nil {
		return nil, err
	}

	config.Token = string(token)

	return &config, nil
}
