package logger

import (
	"context"
	"log/slog"
)

type ContextHandler struct {
	Handler slog.Handler
}

var _ slog.Handler = &ContextHandler{}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *ContextHandler) Handle(ctx context.Context, record slog.Record) error {
	loggerContext := GetLoggerContext(ctx)
	record.AddAttrs(loggerContext.Attrs...)
	return h.Handler.Handle(ctx, record)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{
		Handler: h.Handler.WithAttrs(attrs),
	}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{
		Handler: h.Handler.WithGroup(name),
	}
}
