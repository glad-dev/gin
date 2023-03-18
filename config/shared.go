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

	"golang.org/x/term"
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
		fmt.Printf("%d) %s\n", i+1, config.URL.String())
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

	index-- // Since we add one when displaying the list
	if index < 0 || index >= len(configs) {
		return -1, errors.New("invalid index passed")
	}

	return index, nil
}

func readConfigFromStdIn() (*GitLab, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("What is the base URL (e.g. https://gitlab.com)? ")
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	gitLabURL, err := checkURLStr(scanner.Text())
	if err != nil {
		return nil, err
	}

	// Get the hostname for the API token's name
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}

	fullURL := gitLabURL.JoinPath("-/profile/personal_access_tokens")
	rest := fmt.Sprintf("?name=%s-git-navigator&scopes=api,read_api,read_user", hostname) // Can't be added with url.JoinPath since that escapes the '?'
	fmt.Printf("Go to %s%s to create an API key with the permissions api, read_api and read_user\n", fullURL.String(), rest)

	fmt.Printf("Enter the API token (input is hidden): ")
	token, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Printf("\n")
	if err != nil {
		return nil, err
	}

	return &GitLab{
		URL:   *gitLabURL,
		Token: string(token),
	}, nil
}
