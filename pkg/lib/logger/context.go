package logger

import (
	"context"
	"log/slog"
)

type loggerContextKeyType struct{}

var loggerContextKey = loggerContextKeyType{}

type LoggerContext struct {
	Attrs []slog.Attr
}

func ContextWithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	loggerContext, ok := ctx.Value(loggerContextKey).(*LoggerContext)
	if !ok || loggerContext == nil {
		loggerContext = &LoggerContext{}
	}
	loggerContext.Attrs = append(loggerContext.Attrs, attrs...)
	return context.WithValue(ctx, loggerContextKey, loggerContext)
}

func GetLoggerContext(ctx context.Context) *LoggerContext {
	loggerContext, ok := ctx.Value(loggerContextKey).(*LoggerContext)
	if !ok || loggerContext == nil {
		return &LoggerContext{}
	}
	return loggerContext
}
