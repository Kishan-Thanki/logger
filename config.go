package logger

import (
	"io"
	"log/slog"
)

type Config struct {
	Level      *slog.LevelVar
	RedactKeys []string
	TraceID    bool
	Source     bool
	Output     io.Writer
}

type Option func(*Config)

func WithLevel(level string) Option {
	return func(c *Config) {
		if c.Level == nil {
			c.Level = new(slog.LevelVar)
		}
		var l slog.Level
		if err := l.UnmarshalText([]byte(level)); err == nil {
			c.Level.Set(l)
		}
	}
}

func WithDynamicLevel(lvl *slog.LevelVar) Option {
	return func(c *Config) {
		c.Level = lvl
	}
}

func WithRedaction(keys ...string) Option {
	return func(c *Config) {
		c.RedactKeys = append(c.RedactKeys, keys...)
	}
}

func WithTraceID(enabled bool) Option {
	return func(c *Config) {
		c.TraceID = enabled
	}
}

func WithSource(enabled bool) Option {
	return func(c *Config) {
		c.Source = enabled
	}
}

func WithOutput(w io.Writer) Option {
	return func(c *Config) {
		c.Output = w
	}
}
