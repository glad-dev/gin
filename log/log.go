package log

import (
	"fmt"
	l "log"
	"log/slog"
	"os"
	"path"
	"runtime"

	"github.com/glad-dev/gin/configuration/location"
)

var logger *slog.Logger

// init creates the log directory and file and initializes Log to use said file.
func init() {
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

	logger = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}))
}

func getSource() []any {
	_, file, lineNumber, ok := runtime.Caller(2)
	if ok {
		return []any{"source", fmt.Sprintf("%s:%d", file, lineNumber)}
	}

	return nil
}

func Error(msg string, args ...any) {
	args = append(getSource(), args...)
	logger.Error(msg, args...)
}

func Info(msg string, args ...any) {
	args = append(getSource(), args...)
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	args = append(getSource(), args...)
	logger.Warn(msg, args...)
}
