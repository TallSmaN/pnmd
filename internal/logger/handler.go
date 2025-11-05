package logger

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

type TreeHandler struct {
	cfg Config
}

// NewTreeHandler constructs a TreeHandler using the provided cfg.
// The cfg is normalized to fill any missing fields.
func NewTreeHandler(cfg Config) *TreeHandler {
	return &TreeHandler{cfg: Normalize(cfg)}
}

// Enabled reports whether the given level is enabled by cfg.Level.
func (h *TreeHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.cfg.Level
}

// Handle formats the record as a tree and writes it to stdout.
// Output includes timestamp, level, message, attrs, and optional caller info.
// The context is ignored.
func (h *TreeHandler) Handle(_ context.Context, r slog.Record) error {
	var builder strings.Builder

	builder.WriteString(pterm.Gray(time.Now().Format("2006-01-02 15:04:05")))
	builder.WriteString(" ")

	style := styleForLevel(r.Level)
	builder.WriteString(style.Sprintf("%s", r.Level.String()))
	builder.WriteString(" ")

	builder.WriteString(r.Message)

	var arguments []string
	r.Attrs(func(attr slog.Attr) bool {
		arguments = append(arguments, style.Sprintf("%s: ", attr.Key)+fmt.Sprint(attr.Value))
		return true
	})

	if h.cfg.callerOn(r.Level) {
		if fn := runtime.FuncForPC(r.PC); fn != nil {
			file, line := fn.FileLine(r.PC)
			arguments = append(arguments, pterm.NewStyle(pterm.Bold, pterm.FgGray).
				Sprintf("caller: %s:%d", file, line))
		}
	}

	padding := len(time.Time{}.Format("2006/01/02 15:04:05")) + 3
	for i, arg := range arguments {
		pipe := "└"
		if i < len(arguments)-1 {
			pipe = "├"
		}
		builder.WriteString("\n")
		builder.WriteString(strings.Repeat(" ", padding))
		builder.WriteString(pipe)
		builder.WriteString(" ")
		builder.WriteString(arg)
	}

	fmt.Println(builder.String())
	return nil
}

// WithAttrs returns h unchanged. The handler does not support attribute scoping.
func (h *TreeHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }

// WithGroup returns h unchanged. The handler does not support grouping.
func (h *TreeHandler) WithGroup(_ string) slog.Handler { return h }
