package logging

import (
	"log/slog"
	"os"
)

func New(level slog.Level) *slog.Logger {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}

func Default() *slog.Logger {
	return New(slog.LevelInfo)
}
