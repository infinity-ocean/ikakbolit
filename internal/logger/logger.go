package logger

import (
	"io"
	"log/slog"
	"os"
	"strconv"
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
	level, err := strconv.Atoi(os.Getenv("LOGGING_LEVEL"))
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(mw, &slog.HandlerOptions{
		Level: slog.Level(level),
	})

	return slog.New(handler), nil
}

func (l *logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *logger) Close() error {
	return l.file.Close()
}