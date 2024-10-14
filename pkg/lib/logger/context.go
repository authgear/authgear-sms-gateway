package logger

import (
	"log/slog"
)

type LoggerContexter interface {
	GetAttrs() []slog.Attr
}
