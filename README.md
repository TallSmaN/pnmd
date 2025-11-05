# pnmd
🪵 **My personal slog.Handler with struct fields**

`pnmd` provides a minimal [slog.Handler](https://pkg.go.dev/log/slog#Handler) based on [pterm](https://github.com/pterm/pterm) that prints logs as a compact tree:
timestamp, level, message, structured attributes, and optional caller info.

## Example output
![demo](./.github/assets/output.png)

---

## Example
```go
package main

import (
	"log/slog"

	"github.com/TallSmaN/pnmd"
)

func main() {
	pnmd.Configure(pnmd.Options{
		Level: slog.LevelDebug,
		CallerEnabled: map[slog.Level]bool{
			slog.LevelDebug: true,
			slog.LevelInfo:  true,
			slog.LevelWarn:  true,
			slog.LevelError: true,
		},
	})

	log := pnmd.Get()

	log.Info("server started", "addr", ":8080")
	log.Debug("cache primed", "keys", "123")

	pnmd.SetLevel(slog.LevelWarn)
	log.Warn("slow request", "latency", "230ms")
	log.Error("write failed", "path", "/tmp/out.log")
}
```
---

This project is licensed under the [MIT License](LICENSE).
