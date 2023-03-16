package config

import (
	"errors"
	"fmt"
	"net/url"
	"os/user"
	"path"
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
