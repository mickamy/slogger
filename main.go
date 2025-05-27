package slogger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// Level represents the logging level.
type Level slog.Level

const (
	// LevelDebug represents debug level logs.
	LevelDebug = Level(slog.LevelDebug)

	// LevelInfo represents info level logs.
	LevelInfo = Level(slog.LevelInfo)

	// LevelWarn represents warning level logs.
	LevelWarn = Level(slog.LevelWarn)

	// LevelError represents error level logs.
	LevelError = Level(slog.LevelError)
)

// Config defines the configuration for the slogger library.
type Config struct {
	// Level defines the minimum logging level.
	// The default is LevelInfo.
	Level

	// AddSource indicates whether to add the source file and line number to the log entries.
	// The default is false.
	AddSource bool

	// TrimPathPrefix is the prefix to be trimmed from the file path in the log output.
	TrimPathPrefix string

	// Outputs is the list of writers to which the logs will be written.
	Outputs []io.Writer

	// ContextFieldsExtractor is a function that extracts additional fields from the context.
	ContextFieldsExtractor func(ctx context.Context) []any

	// Handler is the slog.Handler used for logging.
	// If nil, a default JSON handler will be created.
	Handler slog.Handler
}

var (
	// logger is the global instance of slog.Logger.
	logger *slog.Logger

	// cfg is the global configuration for the logger.
	cfg Config
)

func init() {
	SetConfig(ensureDefaults(Config{}))
}

// SetConfig initializes the slogger library with the given configuration.
// This function must be called before using the logger.
func SetConfig(c Config) {
	cfg = ensureDefaults(c)
	logger = slog.New(cfg.Handler)
}

// ensureDefaults ensures the configuration has default values for unset fields.
func ensureDefaults(c Config) Config {
	if c.Outputs == nil {
		c.Outputs = []io.Writer{os.Stdout}
	}
	if c.ContextFieldsExtractor == nil {
		c.ContextFieldsExtractor = func(context.Context) []any { return nil }
	}
	if c.Handler == nil {
		c.Handler = DefaultHandler(c.Level, c.Outputs...)
	}
	return c
}

// DefaultHandler creates a new slog.JSONHandler with the configured options.
func DefaultHandler(level Level, outputs ...io.Writer) *slog.JSONHandler {
	options := &slog.HandlerOptions{
		Level: slog.Level(level),
	}
	writer := io.MultiWriter(outputs...)
	return slog.NewJSONHandler(writer, options)
}

// handle logs a message with the specified level and arguments.
func handle(level Level, msg string, fields ...any) {
	if cfg.AddSource {
		_, file, line, _ := runtime.Caller(2)
		source := fmt.Sprintf("%s:%d", file, line)

		source = strings.TrimPrefix(source, path.Join(cfg.TrimPathPrefix, ""))
		fields = append(fields, slog.String("source", source))
	}

	logger.Log(context.Background(), slog.Level(level), msg, fields...)
}

// handleCtx logs a message with the specified level, arguments, and context.
func handleCtx(ctx context.Context, level Level, msg string, fields ...any) {
	if cfg.AddSource {
		_, file, line, _ := runtime.Caller(2)
		source := fmt.Sprintf("%s:%d", file, line)

		fields = append(fields, cfg.ContextFieldsExtractor(ctx)...)

		source = strings.TrimPrefix(source, filepath.Join(cfg.TrimPathPrefix, ""))
		fields = append(fields, slog.String("source", source))
	}

	logger.Log(ctx, slog.Level(level), msg, fields...)
}

// Debug logs a debug level message.
func Debug(msg string, fields ...any) {
	handle(LevelDebug, msg, fields...)
}

// Info logs an info level message.
func Info(msg string, fields ...any) {
	handle(LevelInfo, msg, fields...)
}

// Warn logs a warning level message.
func Warn(msg string, fields ...any) {
	handle(LevelWarn, msg, fields...)
}

// Error logs an error level message.
func Error(msg string, fields ...any) {
	handle(LevelError, msg, fields...)
}

// DebugCtx logs a debug level message with context.
func DebugCtx(ctx context.Context, msg string, fields ...any) {
	handleCtx(ctx, LevelDebug, msg, fields...)
}

// InfoCtx logs an info level message with context.
func InfoCtx(ctx context.Context, msg string, fields ...any) {
	handleCtx(ctx, LevelInfo, msg, fields...)
}

// WarnCtx logs a warning level message with context.
func WarnCtx(ctx context.Context, msg string, fields ...any) {
	handleCtx(ctx, LevelWarn, msg, fields...)
}

// ErrorCtx logs an error level message with context.
func ErrorCtx(ctx context.Context, msg string, fields ...any) {
	handleCtx(ctx, LevelError, msg, fields...)
}

// StandardLogger creates a standard log.Logger instance with the specified level.
// This is useful for compatibility with libraries expecting a standard logger.
func StandardLogger(level Level, outputs ...io.Writer) *log.Logger {
	return slog.NewLogLogger(DefaultHandler(level, outputs...), slog.Level(level))
}
