package location

import (
	"log" // We can't use out logger since that would lead to an import cycle
	"os"
	"os/user"
	"path"
)

// Dir returns the path to the configuration directory.
func Dir() (string, error) {
	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get the user's home directory: %s", err)
	}

	return path.Join(usr.HomeDir, ".config", "gin"), nil
}

// CreateDir creates the configuration directory.
func CreateDir() error {
	dir, err := Dir()
	if err != nil {
		return err
	}

	return os.MkdirAll(dir, 0o700)
}

// Get returns the path to the configuration file.
func Get() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}

	return path.Join(dir, "config.toml"), nil
}
