package main

import (
	"log/slog"
	"sync"

	log "github.com/TallSmaN/pnmd/internal/logger"
)

type Options = log.Config

type Logger struct {
	mu  sync.RWMutex
	cfg log.Config
	*slog.Logger
}

// defaultLogger is the package-level singleton instance.
var defaultLogger = &Logger{
	cfg: log.DefaultConfig(),
}

// Get returns the global logger.
// Builds it on first use and always returns the same instance.
// Safe for concurrent calls.
func Get() *Logger {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	if defaultLogger.Logger == nil {
		defaultLogger.Logger = slog.New(log.NewTreeHandler(defaultLogger.cfg))
	}
	return defaultLogger
}

// Configure replaces the current configuration and rebuilds the logger.
// Should typically be called once during initialization.
func (l *Logger) Configure(o Options) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cfg = log.Normalize(o)
	l.Logger = slog.New(log.NewTreeHandler(l.cfg))
	return l
}

// SetLevel updates the minimum enabled log level.
// Safe for concurrent use.
func (l *Logger) SetLevel(level slog.Level) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cfg.Level = level
	l.Logger = slog.New(log.NewTreeHandler(l.cfg))
	return l
}

// DisableCallerFor disables caller information output for the specified log levels.
// Automatically initializes the caller map if it is nil.
// Safe for concurrent use.
func (l *Logger) DisableCallerFor(levels ...slog.Level) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ensureCallerMapLocked()

	for _, lv := range levels {
		l.cfg.CallerEnabled[lv] = false
	}

	l.Logger = slog.New(log.NewTreeHandler(l.cfg))
	return l
}

// EnableCallerFor enables caller information output for the specified log levels.
// Automatically initializes the caller map if it is nil.
// Safe for concurrent use.
func (l *Logger) EnableCallerFor(levels ...slog.Level) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ensureCallerMapLocked()

	for _, lv := range levels {
		l.cfg.CallerEnabled[lv] = true
	}

	l.Logger = slog.New(log.NewTreeHandler(l.cfg))
	return l
}

// get returns the internal slog.Logger, creating it if necessary.
func (l *Logger) get() *slog.Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.Logger == nil {
		l.Logger = slog.New(log.NewTreeHandler(l.cfg))
	}
	return l.Logger
}

// ensureCallerMapLocked ensures that the CallerEnabled map exists.
// The caller must hold the write lock.
func (l *Logger) ensureCallerMapLocked() {
	if l.cfg.CallerEnabled == nil {
		l.cfg.CallerEnabled = map[slog.Level]bool{
			slog.LevelDebug: true,
			slog.LevelInfo:  true,
			slog.LevelWarn:  true,
			slog.LevelError: true,
		}
	}
}
