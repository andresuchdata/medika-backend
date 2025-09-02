package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...interface{})
	Info(ctx context.Context, msg string, fields ...interface{})
	Warn(ctx context.Context, msg string, fields ...interface{})
	Error(ctx context.Context, msg string, fields ...interface{})
	Fatal(ctx context.Context, msg string, fields ...interface{})
}

type ZerologLogger struct {
	logger zerolog.Logger
}

func New() Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()

	// Set log level from environment
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		logger = logger.Level(zerolog.DebugLevel)
	case "info":
		logger = logger.Level(zerolog.InfoLevel)
	case "warn":
		logger = logger.Level(zerolog.WarnLevel)
	case "error":
		logger = logger.Level(zerolog.ErrorLevel)
	default:
		logger = logger.Level(zerolog.InfoLevel)
	}

	return &ZerologLogger{logger: logger}
}

func (l *ZerologLogger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	l.logWithTrace(ctx, l.logger.Debug(), msg, fields...)
}

func (l *ZerologLogger) Info(ctx context.Context, msg string, fields ...interface{}) {
	l.logWithTrace(ctx, l.logger.Info(), msg, fields...)
}

func (l *ZerologLogger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	l.logWithTrace(ctx, l.logger.Warn(), msg, fields...)
}

func (l *ZerologLogger) Error(ctx context.Context, msg string, fields ...interface{}) {
	l.logWithTrace(ctx, l.logger.Error(), msg, fields...)
}

func (l *ZerologLogger) Fatal(ctx context.Context, msg string, fields ...interface{}) {
	l.logWithTrace(ctx, l.logger.Fatal(), msg, fields...)
}

func (l *ZerologLogger) logWithTrace(ctx context.Context, event *zerolog.Event, msg string, fields ...interface{}) {
	// Add tracing information if available
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		event = event.
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("span_id", span.SpanContext().SpanID().String())
	}

	// Add fields
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok1 := fields[i].(string)
			value := fields[i+1]
			if ok1 {
				event = event.Interface(key, value)
			}
		}
	}

	event.Msg(msg)
}
