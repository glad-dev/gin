package location

import (
	"fmt"
	"log" // We can't use out logger since that would lead to an import cycle
	"os"
	"os/user"
	"path"
)

func Dir() (string, error) {
	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get the user's home directory: %s", err)

		return "", fmt.Errorf("could not get current user: %w", err)
	}

	return path.Join(usr.HomeDir, ".config", "gn"), nil
}

func CreateDir() error {
	dir, err := Dir()
	if err != nil {
		return err
	}

	return os.MkdirAll(dir, 0o700)
}

func Get() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}

	return path.Join(dir, "gn.toml"), nil
}
