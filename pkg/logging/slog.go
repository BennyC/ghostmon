package logging

import (
	"golang.org/x/exp/slog"
	"io"
	"os"
)

// New will return a new *slog.Logger
func New() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

// NewNilLogger will return a new *slog.Logger that discards all logs wrote
// to the *slog.Logger
func NewNilLogger() *slog.Logger {
	handler := slog.NewTextHandler(io.Discard)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
