package logger

import (
	"context"
	"log"
	"os"
)

const CorrelationIDKey = "correlationID"

type Logger struct {
	*log.Logger
}

func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	correlationID := ctx.Value(CorrelationIDKey)
	if correlationID == nil {
		correlationID = "unknown"
	}
	l.Printf("[%s] INFO: %s %v", correlationID, msg, args)
}

func (l *Logger) Error(ctx context.Context, msg string, err error) {
	correlationID := ctx.Value(CorrelationIDKey)
	if correlationID == nil {
		correlationID = "unknown"
	}
	l.Printf("[%s] ERROR: %s - %v", correlationID, msg, err)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	correlationID := ctx.Value(CorrelationIDKey)
	if correlationID == nil {
		correlationID = "unknown"
	}
	l.Printf("[%s] WARN: %s %v", correlationID, msg, args)
}

func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, correlationID)
}

func GetCorrelationID(ctx context.Context) string {
	id := ctx.Value(CorrelationIDKey)
	if id == nil {
		return "unknown"
	}
	return id.(string)
}
