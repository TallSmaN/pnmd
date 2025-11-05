package pnmd

import (
	"log/slog"

	log "github.com/TallSmaN/pnmd/internal/logger"
)

type Options = log.Config

var (
	cfg    = log.DefaultConfig()
	global *slog.Logger
)

// Get returns the package-level logger, constructing it on first use.
// Not concurrency-safe; call from init or a single goroutine during setup.
func Get() *slog.Logger {
	if global == nil {
		global = slog.New(log.NewTreeHandler(cfg))
	}

	return global
}

// Configure replaces the global configuration and logger.
// Not concurrency-safe; prefer calling during process initialization.
func Configure(o Options) {
	cfg = log.Normalize(o)
	global = slog.New(log.NewTreeHandler(cfg))
}

// SetLevel updates the minimum enabled level on the global logger.
// Not concurrency-safe; racing calls may produce undefined results.
func SetLevel(level slog.Level) {
	cfg.Level = level
	global = slog.New(log.NewTreeHandler(cfg))
}

// DisableCallerFor disables caller display for the given levels.
// If the caller map is nil, defaults are restored first.
// Not concurrency-safe.
func DisableCallerFor(levels ...slog.Level) {
	if cfg.CallerEnabled == nil {
		cfg = log.DefaultConfig()
	}

	for _, lv := range levels {
		cfg.CallerEnabled[lv] = false
	}

	global = slog.New(log.NewTreeHandler(cfg))
}

// EnableCallerFor enables caller display for the given levels.
// If the caller map is nil, defaults are restored first.
// Not concurrency-safe.
func EnableCallerFor(levels ...slog.Level) {
	if cfg.CallerEnabled == nil {
		cfg = log.DefaultConfig()
	}

	for _, lv := range levels {
		cfg.CallerEnabled[lv] = true
	}

	global = slog.New(log.NewTreeHandler(cfg))
}
