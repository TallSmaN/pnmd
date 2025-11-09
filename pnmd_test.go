package pnmd_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/TallSmaN/pnmd"
)

func TestExample(t *testing.T) {
	logger := slog.New(
		pnmd.NewHandler(os.Stdout, &pnmd.Options{
			Level: slog.LevelDebug,
			CallerEnabled: map[slog.Level]bool{
				slog.LevelDebug: true,
				slog.LevelInfo:  true,
				slog.LevelWarn:  true,
				slog.LevelError: true,
			},
			TimeFormat: "2006/01/02 15:04:05",
			Padding:    3,
		}),
	)

	logger.Debug("initializing cache subsystem", "cache", "redis", "host", "localhost", "port", 6379) //nolint

	logger.Info("cache connected") //nolint

	logger.Warn("slow query detected", "duration_ms", 1823, "query", "SELECT * FROM users WHERE active=1", "user", "analytics-worker") //nolint

	logger.Error("failed to write audit event", "error", "disk full", "path", "/var/log/audit.json", "component", "audit", "retry_in_sec", 30) //nolint

	logger.Debug("config reloaded", "file", "/etc/app/config.yaml", "changes", 5) //nolint

	logger.Info("http server started", "addr", ":8080", "threads", 8) //nolint

	logger.Warn("deprecated API usage", "endpoint", "/v1/legacy", "client", "mobile-android", "version", "1.2.0") //nolint

	logger.Error("user authentication failed", "user", "john", "ip", "192.168.1.42", "reason", "invalid token") //nolint

	logger.Debug("background job finished", "job_id", "import-2025-11-05", "rows", 152_000, "duration_sec", 94.2) //nolint

	logger.Info("graceful shutdown complete", "uptime_min", 238) //nolint

	t.Skip()
}
