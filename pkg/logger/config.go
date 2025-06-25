package logger

import (
	"os"
	"strings"
)

// Environment constants
const (
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
	EnvTest        = "test"
)

// Log level constants
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
	LevelPanic = "panic"
)

// Encoding constants
const (
	EncodingJSON    = "json"
	EncodingConsole = "console"
)

// ConfigFromEnv creates logger configuration from environment variables
func ConfigFromEnv() Config {
	config := DefaultConfig()

	// Get environment
	if env := os.Getenv("APP_ENV"); env != "" {
		config.Environment = strings.ToLower(env)
	}
	if env := os.Getenv("GIN_MODE"); env != "" {
		switch env {
		case "release":
			config.Environment = EnvProduction
		case "test":
			config.Environment = EnvTest
		default:
			config.Environment = EnvDevelopment
		}
	}

	// Get log level
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Level = strings.ToLower(level)
	}

	// Get encoding
	if encoding := os.Getenv("LOG_ENCODING"); encoding != "" {
		config.Encoding = strings.ToLower(encoding)
	}

	// Get output paths
	if outputs := os.Getenv("LOG_OUTPUT_PATHS"); outputs != "" {
		config.OutputPaths = strings.Split(outputs, ",")
	}

	// Adjust config based on environment
	switch config.Environment {
	case EnvProduction:
		if config.Level == "" {
			config.Level = LevelInfo
		}
		config.Encoding = EncodingJSON
	case EnvStaging:
		if config.Level == "" {
			config.Level = LevelInfo
		}
		config.Encoding = EncodingJSON
	case EnvTest:
		if config.Level == "" {
			config.Level = LevelError
		}
		config.Encoding = EncodingConsole
	default: // development
		if config.Level == "" {
			config.Level = LevelDebug
		}
		config.Encoding = EncodingConsole
	}

	return config
}

// DevelopmentConfig returns configuration optimized for development
func DevelopmentConfig() Config {
	return Config{
		Level:       LevelDebug,
		Environment: EnvDevelopment,
		OutputPaths: []string{"stdout"},
		Encoding:    EncodingConsole,
	}
}

// ProductionConfig returns configuration optimized for production
func ProductionConfig() Config {
	return Config{
		Level:       LevelInfo,
		Environment: EnvProduction,
		OutputPaths: []string{"stdout"},
		Encoding:    EncodingJSON,
	}
}

// TestConfig returns configuration optimized for testing
func TestConfig() Config {
	return Config{
		Level:       LevelError,
		Environment: EnvTest,
		OutputPaths: []string{"stdout"},
		Encoding:    EncodingConsole,
	}
}

// IsProduction checks if the environment is production
func (c Config) IsProduction() bool {
	return c.Environment == EnvProduction
}

// IsDevelopment checks if the environment is development
func (c Config) IsDevelopment() bool {
	return c.Environment == EnvDevelopment
}

// IsTest checks if the environment is test
func (c Config) IsTest() bool {
	return c.Environment == EnvTest
}

// Validate validates the configuration
func (c Config) Validate() error {
	// Validate log level
	validLevels := map[string]bool{
		LevelDebug: true,
		LevelInfo:  true,
		LevelWarn:  true,
		LevelError: true,
		LevelFatal: true,
		LevelPanic: true,
	}
	if !validLevels[c.Level] {
		c.Level = LevelInfo
	}

	// Validate environment
	validEnvs := map[string]bool{
		EnvDevelopment: true,
		EnvStaging:     true,
		EnvProduction:  true,
		EnvTest:        true,
	}
	if !validEnvs[c.Environment] {
		c.Environment = EnvDevelopment
	}

	// Validate encoding
	validEncodings := map[string]bool{
		EncodingJSON:    true,
		EncodingConsole: true,
	}
	if !validEncodings[c.Encoding] {
		if c.IsProduction() {
			c.Encoding = EncodingJSON
		} else {
			c.Encoding = EncodingConsole
		}
	}

	// Validate output paths
	if len(c.OutputPaths) == 0 {
		c.OutputPaths = []string{"stdout"}
	}

	return nil
}

// WithLevel sets the log level
func (c Config) WithLevel(level string) Config {
	c.Level = strings.ToLower(level)
	return c
}

// WithEnvironment sets the environment
func (c Config) WithEnvironment(env string) Config {
	c.Environment = strings.ToLower(env)
	return c
}

// WithEncoding sets the encoding
func (c Config) WithEncoding(encoding string) Config {
	c.Encoding = strings.ToLower(encoding)
	return c
}

// WithOutputPaths sets the output paths
func (c Config) WithOutputPaths(paths ...string) Config {
	c.OutputPaths = paths
	return c
}

// AddOutputPath adds an output path
func (c Config) AddOutputPath(path string) Config {
	c.OutputPaths = append(c.OutputPaths, path)
	return c
}

// GetEffectiveConfig returns the effective configuration after applying defaults and validation
func GetEffectiveConfig() Config {
	config := ConfigFromEnv()
	config.Validate()
	return config
}
