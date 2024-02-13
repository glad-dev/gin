package location

import (
	"log" // We can't use our log since that would lead to an import cycle
	"os"
	"os/user"
	"path"
)

// Dir returns the path to the configuration directory. If $XDG_CONFIG_HOME is set, $XDG_CONFIG_HOME/gin is used.
// Otherwise, $HOME/.config/gin is used.
func Dir() (string, error) {
	// Check if $XDG_CONFIG_HOME is set
	val, ok := os.LookupEnv("XDG_CONFIG_HOME")
	if ok && len(val) > 0 {
		return path.Join(val, "gin"), nil
	}

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
