package logger

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	stderrHandler := slog.NewTextHandler(os.Stderr, nil)
	handler := &ContextHandler{
		Handler: stderrHandler,
	}
	return slog.New(handler)
}
