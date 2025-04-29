package logger

import (
	"io"
	"log/slog"
	"os"
)

type logger struct {
	logger *slog.Logger
	file   *os.File
}

func NewLogger() (*slog.Logger, error) {
	file, err := os.OpenFile(
		os.Getenv("LOGGING_FILE_PATH"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
        0o644)
	if err != nil {
		return nil, err
	}

	mw := io.MultiWriter(os.Stdout, file)

	var level slog.Level
	levelStr := os.Getenv("LOGGING_LEVEL")

	switch levelStr {
	case "ERROR":
		level = slog.LevelError
	case "DEBUG":
		level = slog.LevelDebug
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(mw, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler), nil
}