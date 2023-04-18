package logger

import (
	l "log"
	"os"
	"path"
	"time"

	"github.com/charmbracelet/log"
	"gn/config/location"
)

var Log *log.Logger

func Init() {
	err := location.CreateDir()
	if err != nil {
		l.Fatalf("Failed to create config directory: %s", err)
	}

	dir, err := location.Dir()
	if err != nil {
		l.Fatalf("Failed to resolve config directory: %s", err)
	}

	file, err := os.OpenFile(path.Join(dir, "log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		l.Fatalf("Failed to open log file: %s", err)
	}

	Log = log.NewWithOptions(file, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.RFC3339,
	})
}
