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
	// LogLevel defines the minimum logging level.
	Level

	// Output is the destination for log output.
	// If not set, os.Stdout is used.
	Output io.Writer

	// TrimPathPrefix is the prefix to be trimmed from the file path in the log output.
	TrimPathPrefix string

	// ContextFieldsExtractor is a function that extracts additional fields from the context.
	ContextFieldsExtractor func(ctx context.Context) []any
}

var (
	// logger is the global instance of slog.Logger.
	logger *slog.Logger

	// config is the global configuration for the logger.
	config Config
)

func init() {
	Init(Config{})
}

// Init initializes the slogger library with the given configuration.
// This function must be called before using the logger.
func Init(cfg Config) {
	config = ensureDefaults(cfg)
	logger = slog.New(createHandler())
}

// ensureDefaults ensures the configuration has default values for unset fields.
func ensureDefaults(cfg Config) Config {
	if cfg.Level == 0 {
		cfg.Level = LevelInfo
	}
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}
	if cfg.ContextFieldsExtractor == nil {
		cfg.ContextFieldsExtractor = func(context.Context) []any { return nil }
	}
	return cfg
}

// createHandler creates a new slog.JSONHandler with the configured options.
func createHandler() *slog.JSONHandler {
	options := &slog.HandlerOptions{
		Level: slog.Level(config.Level),
	}
	writer := io.MultiWriter(config.Output, os.Stdout)
	return slog.NewJSONHandler(writer, options)
}

// handle logs a message with the specified level and arguments.
func handle(level Level, msg string, fields ...any) {
	_, file, line, _ := runtime.Caller(2)
	source := fmt.Sprintf("%s:%d", file, line)

	source = strings.TrimPrefix(source, path.Join(config.TrimPathPrefix, ""))
	fields = append(fields, slog.String("source", source))

	logger.Log(context.Background(), slog.Level(level), msg, fields...)
}

// handleCtx logs a message with the specified level, arguments, and context.
func handleCtx(ctx context.Context, level Level, msg string, fields ...any) {
	_, file, line, _ := runtime.Caller(2)
	source := fmt.Sprintf("%s:%d", file, line)

	fields = append(fields, config.ContextFieldsExtractor(ctx)...)

	source = strings.TrimPrefix(source, filepath.Join(config.TrimPathPrefix, ""))
	fields = append(fields, slog.String("source", source))

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
func StandardLogger(level Level) *log.Logger {
	return slog.NewLogLogger(createHandler(), slog.Level(level))
}
