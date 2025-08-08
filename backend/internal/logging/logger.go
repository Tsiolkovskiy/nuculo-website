package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/google/uuid"
)

// LogLevel represents different logging levels
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
	LevelFatal LogLevel = "FATAL"
)

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
	service string
}

// Config holds logger configuration
type Config struct {
	Level       LogLevel
	Service     string
	Environment string
	Format      string // "json" or "text"
}

// NewLogger creates a new structured logger
func NewLogger(config Config) *Logger {
	var level slog.Level
	switch config.Level {
	case LevelDebug:
		level = slog.LevelDebug
	case LevelInfo:
		level = slog.LevelInfo
	case LevelWarn:
		level = slog.LevelWarn
	case LevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create handler options
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize timestamp format
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   "timestamp",
					Value: slog.StringValue(time.Now().UTC().Format(time.RFC3339)),
				}
			}
			// Customize source format
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				return slog.Attr{
					Key:   "source",
					Value: slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line)),
				}
			}
			return a
		},
	}

	// Create handler based on format
	var handler slog.Handler
	if config.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	// Create logger with service context
	logger := slog.New(handler).With(
		"service", config.Service,
		"environment", config.Environment,
	)

	return &Logger{
		Logger:  logger,
		service: config.Service,
	}
}

// WithContext adds context information to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	attrs := []slog.Attr{}

	// Add request ID if available
	if requestID := GetRequestID(ctx); requestID != "" {
		attrs = append(attrs, slog.String("request_id", requestID))
	}

	// Add user ID if available
	if userID := GetUserID(ctx); userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}

	// Add operation name if available
	if operation := GetOperationName(ctx); operation != "" {
		attrs = append(attrs, slog.String("operation", operation))
	}

	if len(attrs) > 0 {
		return &Logger{
			Logger:  l.Logger.With(attrs...),
			service: l.service,
		}
	}

	return l
}

// WithFields adds structured fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	attrs := make([]slog.Attr, 0, len(fields))
	for key, value := range fields {
		attrs = append(attrs, slog.Any(key, value))
	}

	return &Logger{
		Logger:  l.Logger.With(attrs...),
		service: l.service,
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(msg, args...))
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.Logger.Info(fmt.Sprintf(msg, args...))
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(msg, args...))
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, args...))
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf("FATAL: "+msg, args...))
	os.Exit(1)
}

// LogError logs an error with stack trace
func (l *Logger) LogError(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(1)
	if ok {
		l.Logger.Error(
			fmt.Sprintf(msg, args...),
			"error", err.Error(),
			"caller", fmt.Sprintf("%s:%d", file, line),
		)
	} else {
		l.Logger.Error(
			fmt.Sprintf(msg, args...),
			"error", err.Error(),
		)
	}
}

// Context key types
type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
	operationKey contextKey = "operation_name"
	loggerKey    contextKey = "logger"
)

// Context helpers
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(userIDKey).(string); ok {
		return id
	}
	return ""
}

func WithOperationName(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, operationKey, operation)
}

func GetOperationName(ctx context.Context) string {
	if op, ok := ctx.Value(operationKey).(string); ok {
		return op
	}
	return ""
}

func WithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerKey).(*Logger); ok {
		return logger
	}
	// Return default logger if none in context
	return NewLogger(Config{
		Level:       LevelInfo,
		Service:     "unknown",
		Environment: "development",
		Format:      "text",
	})
}

// GenerateRequestID generates a new request ID
func GenerateRequestID() string {
	return uuid.New().String()
}