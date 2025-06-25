package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger interface for dependency injection
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sync() error
}

// ZapLogger wraps zap.Logger to implement our Logger interface
type ZapLogger struct {
	logger *zap.Logger
}

// Global logger instance
var globalLogger Logger

// Config holds logger configuration
type Config struct {
	Level       string   `json:"level" yaml:"level"`
	Environment string   `json:"environment" yaml:"environment"`
	OutputPaths []string `json:"output_paths" yaml:"output_paths"`
	Encoding    string   `json:"encoding" yaml:"encoding"`
}

// DefaultConfig returns default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:       "info",
		Environment: "development",
		OutputPaths: []string{"stdout"},
		Encoding:    "console",
	}
}

// Initialize initializes the global logger with the given configuration
func Initialize(config Config) error {
	logger, err := NewLogger(config)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// NewLogger creates a new logger instance with the given configuration
func NewLogger(config Config) (Logger, error) {
	// Parse log level
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// Create encoder config based on environment
	var encoderConfig zapcore.EncoderConfig
	if config.Environment == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()
		config.Encoding = "json"
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Configure time encoding
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create encoder
	var encoder zapcore.Encoder
	if config.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Create writer syncer
	var writeSyncer zapcore.WriteSyncer
	if len(config.OutputPaths) == 0 || (len(config.OutputPaths) == 1 && config.OutputPaths[0] == "stdout") {
		writeSyncer = zapcore.AddSync(os.Stdout)
	} else {
		// For file outputs, you might want to add file rotation here
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// Create logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{logger: zapLogger}, nil
}

// Implementation of Logger interface

func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *ZapLogger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

func (l *ZapLogger) With(fields ...zap.Field) Logger {
	return &ZapLogger{logger: l.logger.With(fields...)}
}

func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}

// Global logger functions

// GetLogger returns the global logger instance
func GetLogger() Logger {
	if globalLogger == nil {
		// Initialize with default config if not initialized
		config := DefaultConfig()
		if env := os.Getenv("APP_ENV"); env != "" {
			config.Environment = env
		}
		if level := os.Getenv("LOG_LEVEL"); level != "" {
			config.Level = strings.ToLower(level)
		}
		Initialize(config)
	}
	return globalLogger
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Panic logs a panic message and panics
func Panic(msg string, fields ...zap.Field) {
	GetLogger().Panic(msg, fields...)
}

// With creates a child logger with additional fields
func With(fields ...zap.Field) Logger {
	return GetLogger().With(fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return GetLogger().Sync()
}

// Common field helpers

// String creates a string field
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int creates an int field
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Int64 creates an int64 field
func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

// Uint creates a uint field
func Uint(key string, val uint) zap.Field {
	return zap.Uint(key, val)
}

// Uint32 creates a uint32 field
func Uint32(key string, val uint32) zap.Field {
	return zap.Uint32(key, val)
}

// Uint64 creates a uint64 field
func Uint64(key string, val uint64) zap.Field {
	return zap.Uint64(key, val)
}

// Float64 creates a float64 field
func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

// Bool creates a bool field
func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

// Any creates a field with any value
func Any(key string, val any) zap.Field {
	return zap.Any(key, val)
}

// Error creates an error field
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Duration creates a duration field
func Duration(key string, val any) zap.Field {
	return zap.Any(key, val)
}
