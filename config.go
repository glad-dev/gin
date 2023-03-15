package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"golang.org/x/term"
)

type GitlabConfig struct {
	Url          string
	Token        string
	MajorVersion int
}

const currentMajorVersion = 1

var (
	ErrConfigDoesNotExist = errors.New("config does not exist")
	ErrDontCreateConfig   = errors.New("user does not want to create the configuration")
)

func (config *GitlabConfig) CheckValidity() error {
	// Check URL
	_, err := checkURLStr(config.Url)
	if err != nil {
		return err
	}

	// Check if token is semantically correct. The tokens validity is not checked
	if len(config.Token) < 20 {
		return fmt.Errorf("config contains token that is too short. Expected: at least 20, got %d", len(config.Token))
	}

	// Check version
	if config.MajorVersion > currentMajorVersion {
		return fmt.Errorf("config was written by a newer version of the tool")
	}

	return nil
}

func handleConfig() (*GitlabConfig, error) {
	config, err := loadConfig()
	if err != nil {
		if !errors.Is(err, ErrConfigDoesNotExist) {
			// Error is NOT about the config not existing
			return nil, err
		}

		config, err = readConfigFromStdin()
		if err != nil {
			if errors.Is(err, ErrDontCreateConfig) {
				os.Exit(0)
			}

			return nil, err
		}

		err = writeConfig(config)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func writeConfig(config *GitlabConfig) error {
	config.MajorVersion = currentMajorVersion

	err := config.CheckValidity()
	if err != nil {
		return fmt.Errorf("passed config is invalid: %w", err)
	}

	buf := new(bytes.Buffer)
	err = toml.NewEncoder(buf).Encode(config)
	if err != nil {
		return fmt.Errorf("could not encode config: %w", err)
	}

	fileLocation, err := getConfigLocation()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("could not open config file: %w", err)
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}

func loadConfig() (*GitlabConfig, error) {
	fileLocation, err := getConfigLocation()
	if err != nil {
		return nil, err
	}

	// Load config
	config := GitlabConfig{}
	metaData, err := toml.DecodeFile(fileLocation, &config)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrConfigDoesNotExist
		}

		return nil, fmt.Errorf("could not decode config: %w", err)
	}

	// Check if the config only contains the keys we expect
	if len(metaData.Undecoded()) > 0 {
		return nil, fmt.Errorf("config contains unexpected keys: %+v", metaData.Undecoded())
	}

	// Check if the config contains all the keys we need
	expected := reflect.ValueOf(config).NumField()
	if len(metaData.Keys()) != expected {
		return nil, fmt.Errorf("config contains an invalid amount of keys. Expect %d, got %d", expected, len(metaData.Keys()))
	}

	err = config.CheckValidity()
	if err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	return &config, nil
}

func readConfigFromStdin() (*GitlabConfig, error) {
	path, err := getConfigLocation()
	if err != nil {
		return nil, err
	}

	fmt.Printf("No configuration file was found at: %s\n", path)
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

	config := GitlabConfig{
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

func checkURLStr(urlStr string) (*url.URL, error) {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, err
	}

	if !u.IsAbs() {
		return nil, errors.New("URL is not absolute")
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return nil, errors.New("URL has invalid scheme")
	}

	return u, nil
}

func getConfigLocation() (string, error) {
	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("could not get current user: %w", err)
	}

	return path.Join(usr.HomeDir, ".gn.toml"), nil
}
