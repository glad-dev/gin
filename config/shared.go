package config

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path"
	"strconv"
)

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

func selectExistingConfigs(configs []GitLab) (int, error) {
	fmt.Println("The following URLs exist:")
	for i, config := range configs {
		fmt.Printf("%d) %s\n", i+1, config.Url)
	}

	fmt.Print("Enter the index of the config you wish to edit: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return -1, err
	}

	index, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return -1, err
	}

	index -= 1 // Since we add one when displaying the list
	if index < 0 || index >= len(configs) {
		return -1, errors.New("invalid index passed")
	}

	return index, nil
}
