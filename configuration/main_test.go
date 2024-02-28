package configuration

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %s", err)
	}

	err = os.Setenv("XDG_CONFIG_HOME", dir)
	if err != nil {
		log.Fatalf("Failed to set env: %s", err)
	}

	os.Exit(m.Run())
}
