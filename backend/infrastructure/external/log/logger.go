package log

import (
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

func NewLogger() *slog.Logger {
	logLevelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))
	var logLevel slog.Level

	switch logLevelStr {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}))
}
