package logger

import (
	"log/slog"
	"os"
)

func NewTextHandler() *slog.TextHandler {
	return slog.NewTextHandler(os.Stderr, nil)
}
