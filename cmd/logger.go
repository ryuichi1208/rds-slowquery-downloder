package cmd

import (
	"log/slog"
	"os"
)

func NewLogger(logLevel string) *slog.Logger {
	level := slog.LevelInfo

	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	opts := slog.HandlerOptions{
		Level: level,
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &opts))
}
