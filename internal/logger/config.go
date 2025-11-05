package logger

import "log/slog"

type Config struct {
	Level         slog.Level
	CallerEnabled map[slog.Level]bool
}

// DefaultConfig returns a Config with Level=Info and caller enabled for all standard levels.
func DefaultConfig() Config {
	return Config{
		Level: slog.LevelInfo,
		CallerEnabled: map[slog.Level]bool{
			slog.LevelDebug: true,
			slog.LevelInfo:  true,
			slog.LevelWarn:  true,
			slog.LevelError: true,
		},
	}
}

// Normalize ensures the Config is usable. If CallerEnabled is nil, it is set to defaults.
func Normalize(c Config) Config {
	if c.CallerEnabled == nil {
		c.CallerEnabled = DefaultConfig().CallerEnabled
	}
	return c
}

// callerOn reports whether caller info is enabled for the given level.
// If the level key is absent, it returns true.
func (c Config) callerOn(l slog.Level) bool {
	on, ok := c.CallerEnabled[l]
	if !ok {
		return true
	}
	return on
}
