package log

import "log/slog"

const (
	ModeStdio = "stdio"
)

type slogConfig struct {
	path     string
	logLevel slog.Leveler
}

// GetPath...
func (s slogConfig) GetPath() string {
	return s.path
}

// Config...
type Config struct {
	mode string
	slog slogConfig
}

// GetLogLevel...
func (c Config) GetLogLevel() slog.Leveler {
	return c.slog.logLevel
}

// GetMode...
func (c Config) GetMode() string {
	return c.mode
}

// GetSlogConfig
func (c Config) GetSlogConfig() slogConfig {
	return c.slog
}

type ConfigOption func(*Config)

// WithMode...
func WithMode(mode string) ConfigOption {
	return func(c *Config) {
		c.mode = mode
	}
}

// WithLogPath...
func WithLogPath(path string) ConfigOption {
	return func(c *Config) {
		c.slog.path = path
	}
}

// WithLogLevel...
func WithLogLevel(level slog.Level) ConfigOption {
	return func(c *Config) {
		c.slog.logLevel = level
	}
}

// NewConfig...
func NewConfig(opts ...ConfigOption) *Config {
	config := &Config{
		mode: ModeStdio,
		slog: slogConfig{
			logLevel: slog.LevelInfo,
		},
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}
