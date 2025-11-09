package builder

import (
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/TallSmaN/pnmd/internal/opts"
	"github.com/pterm/pterm"
)

// Builder implements LogBuilder to build formatted log entries.
type Builder struct {
	r        *slog.Record
	sb       *strings.Builder
	spacePad string
	argCount int
	style    pterm.RGBStyle
	opts     *opts.Options
}

// NewBuilder creates a new Builder instance for the given record and options.
func NewBuilder(opts *opts.Options, r *slog.Record) *Builder {
	var sb strings.Builder

	includeCall := opts.CallerEnabled[r.Level]
	argCount := r.NumAttrs()

	if includeCall {
		argCount++
	}

	sb.Grow(64 + len(r.Message) + argCount*32)
	padding := len(time.Time{}.Format(opts.TimeFormat)) + opts.Padding
	spacePad := strings.Repeat(" ", padding)
	style := StyleForLevel(r.Level)

	return &Builder{
		opts:     opts,
		r:        r,
		sb:       &sb,
		spacePad: spacePad,
		argCount: argCount,
		style:    style,
	}
}

// WriteTime appends the formatted log timestamp.
func (b *Builder) WriteTime() {
	b.sb.WriteString(pterm.Gray(b.r.Time.Format(b.opts.TimeFormat)))
	b.sb.WriteByte(' ')
}

// WriteLevel appends the log level label with styling.
func (b *Builder) WriteLevel() {
	lvl := b.r.Level.String()

	if len(lvl) > 4 {
		lvl = lvl[:4]
	}

	b.sb.WriteString(b.style.Sprint(lvl))
	b.sb.WriteByte(' ')
}

// WriteMessage appends the main log message.
func (b *Builder) WriteMessage() {
	b.sb.WriteString(b.r.Message)
}

// WriteAttrs appends log attributes and optional caller information.
func (b *Builder) WriteAttrs() {
	i := 0

	b.r.Attrs(func(attr slog.Attr) bool {
		pipe := "└"
		if i < b.argCount-1 {
			pipe = "├"
		}
		b.sb.WriteByte('\n')
		b.sb.WriteString(b.spacePad)
		b.sb.WriteString(pipe)
		b.sb.WriteByte(' ')
		b.sb.WriteString(b.style.Sprint(attr.Key, ": "))
		b.sb.WriteString(attr.Value.String())
		i++

		return true
	})

	if !b.opts.CallerEnabled[b.r.Level] {
		return
	}

	frame, _ := runtime.CallersFrames([]uintptr{b.r.PC}).Next()
	short := filepath.Base(filepath.Dir(frame.File)) + string(filepath.Separator) + filepath.Base(frame.File)
	b.sb.WriteString("\n" + b.spacePad + "└ " + pterm.NewStyle(pterm.Italic, pterm.FgGray).Sprintf("caller: %s:%d", short, frame.Line))
}

// Build finalizes the log entry and returns it as a string.
func (b *Builder) Build() string {
	b.sb.WriteByte('\n')
	return b.sb.String()
}
