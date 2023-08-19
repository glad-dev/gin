package logger

import (
	l "log"
	"log/slog"
	"os"
	"path"

	"github.com/glad-dev/gin/config/location"
)

// Log is the global charmbracelet logger.
var Log *slog.Logger

// Init creates the log directory and file and initializes Log to use said file.
func Init() {
	err := location.CreateDir()
	if err != nil {
		l.Fatalf("Failed to create config directory: %s", err)
	}

	dir, err := location.Dir()
	if err != nil {
		l.Fatalf("Failed to resolve config directory: %s", err)
	}

	file, err := os.OpenFile(path.Join(dir, "gin.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		l.Fatalf("Failed to open log file: %s", err)
	}

	Log = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))
}
