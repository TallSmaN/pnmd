package logger

import (
	"log/slog"

	"github.com/pterm/pterm"
)

// styleForLevel returns a pterm rgb style for the given slog level.
func styleForLevel(l slog.Level) pterm.RGBStyle {
	switch l {
	case slog.LevelDebug:
		return pterm.NewRGB(95, 214, 254).ToRGBStyle().AddOptions(pterm.Bold)
	case slog.LevelInfo:
		return pterm.NewRGB(95, 254, 135).ToRGBStyle().AddOptions(pterm.Bold)
	case slog.LevelWarn:
		return pterm.NewRGB(254, 241, 95).ToRGBStyle().AddOptions(pterm.Bold)
	case slog.LevelError:
		return pterm.NewRGB(254, 95, 134).ToRGBStyle().AddOptions(pterm.Bold)
	default:
		return pterm.NewRGB(255, 255, 255).ToRGBStyle().AddOptions(pterm.Bold)
	}
}
