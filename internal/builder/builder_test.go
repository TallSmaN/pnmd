package builder

import (
	"regexp"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"log/slog"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/TallSmaN/pnmd/internal/opts"
	"github.com/pterm/pterm"
)

func TestBuilder_Build(t *testing.T) {
	tests := []struct {
		name string
		init string
		want string
	}{
		{
			name: "empty builder",
			init: "",
			want: "\n",
		},
		{
			name: "non-empty builder",
			init: "some log output",
			want: "some log output\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			sb.WriteString(tt.init)

			b := &Builder{sb: &sb}
			got := b.Build()

			if got != tt.want {
				t.Errorf("Build() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuilder_WriteAttrs(t *testing.T) {
	pc := func() uintptr { p, _, _, _ := runtime.Caller(0); return p }()

	tests := []struct {
		name          string
		attrs         []any
		level         slog.Level
		callerEnabled bool
		padding       int
		wantSubstrs   []string
	}{
		{
			name:          "single attribute no caller",
			attrs:         []any{"id", "42"},
			level:         slog.LevelInfo,
			callerEnabled: false,
			padding:       3,
			wantSubstrs:   []string{"└", "id: 42"},
		},
		{
			name:          "two attributes with caller",
			attrs:         []any{"user", "alice", "status", "ok"},
			level:         slog.LevelInfo,
			callerEnabled: true,
			padding:       3,
			wantSubstrs:   []string{"├", "user: alice", "status: ok", "caller:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			r := &slog.Record{Level: tt.level, PC: pc}
			r.Add(tt.attrs...)

			o := &opts.Options{
				CallerEnabled: map[slog.Level]bool{
					tt.level: tt.callerEnabled,
				},
				Padding: tt.padding,
			}

			b := &Builder{
				r:        r,
				sb:       &sb,
				spacePad: strings.Repeat(" ", tt.padding),
				argCount: r.NumAttrs(),
				style:    StyleForLevel(tt.level),
				opts:     o,
			}

			b.WriteAttrs()
			out := removeANSI(b.sb.String())

			for _, want := range tt.wantSubstrs {
				if !strings.Contains(out, want) {
					t.Errorf("WriteAttrs() missing substring %q\noutput:\n%s", want, out)
				}
			}
		})
	}
}

func TestBuilder_WriteLevel(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
		style pterm.RGBStyle
		want  string
	}{
		{
			name:  "info level",
			level: slog.LevelInfo,
			style: pterm.NewRGB(95, 254, 135).ToRGBStyle().AddOptions(pterm.Bold),
			want:  pterm.NewRGB(95, 254, 135).ToRGBStyle().AddOptions(pterm.Bold).Sprint("INFO"),
		},
		{
			name:  "debug level",
			level: slog.LevelDebug,
			style: pterm.NewRGB(135, 135, 255).ToRGBStyle().AddOptions(pterm.Bold),
			want:  pterm.NewRGB(135, 135, 255).ToRGBStyle().AddOptions(pterm.Bold).Sprint("DEBU"),
		},
		{
			name:  "warn level",
			level: slog.LevelWarn,
			style: pterm.NewRGB(255, 255, 0).ToRGBStyle().AddOptions(pterm.Bold),
			want:  pterm.NewRGB(255, 255, 0).ToRGBStyle().AddOptions(pterm.Bold).Sprint("WARN"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			r := &slog.Record{Level: tt.level}

			b := &Builder{
				r:     r,
				sb:    &sb,
				style: tt.style,
			}

			b.WriteLevel()
			got := b.sb.String()

			want := tt.want + " "

			if got != want {
				t.Errorf("WriteLevel() mismatch:\n got:  %q\n want: %q", got, want)
			}
		})
	}
}

func TestBuilder_WriteMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "basic",
			message: "its a test log message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder

			r := &slog.Record{
				Message: tt.message,
			}

			b := &Builder{
				r:  r,
				sb: &sb,
			}

			b.WriteMessage()

			got := b.sb.String()
			if got != tt.message {
				t.Errorf("WriteMessage() mismatch:\n got:  %q\n want: %q", got, tt.message)
			}
		})
	}
}

func TestBuilder_WriteTime(t *testing.T) {
	fakeTime := time.Date(2009, 11, 11, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		timeFormat string
		want       string
	}{
		{
			name:       "default time format",
			timeFormat: "2006/01/02 15:04:05",
			want:       pterm.Gray("2009/11/11 00:00:00"),
		},
		{
			name:       "kitchen time format",
			timeFormat: time.Kitchen,
			want:       pterm.Gray("12:00AM"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder

			r := &slog.Record{
				Time: fakeTime,
			}

			b := &Builder{
				r:    r,
				sb:   &sb,
				opts: &opts.Options{TimeFormat: tt.timeFormat},
			}

			b.WriteTime()
			got := b.sb.String()

			want := tt.want + " "

			if got != want {
				t.Errorf("WriteTime() mismatch:\n got:  %q\n want: %q", got, want)
			}
		})
	}
}

func TestNewBuilder(t *testing.T) {
	pc := func() uintptr { p, _, _, _ := runtime.Caller(0); return p }()
	fakeTime := time.Date(2009, 11, 11, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		opts       *opts.Options
		makeRecord func() *slog.Record
		want       *Builder
		wantCap    int
	}{
		{
			name: "default options",
			opts: opts.Default(),
			makeRecord: func() *slog.Record {
				r := &slog.Record{
					Time:    fakeTime,
					Message: "default log",
					Level:   slog.LevelInfo,
					PC:      pc,
				}
				return r
			},
			wantCap: 64 + len("default log") + 1*32,
		},
		{
			name: "custom options with 2 attrs",
			opts: &opts.Options{
				Level: slog.LevelDebug,
				CallerEnabled: map[slog.Level]bool{
					slog.LevelDebug: true,
					slog.LevelInfo:  true,
					slog.LevelWarn:  true,
					slog.LevelError: true,
				},
				TimeFormat: time.Kitchen,
				Padding:    4,
			},
			makeRecord: func() *slog.Record {
				r := &slog.Record{
					Time:    fakeTime,
					Message: "custom log message",
					Level:   slog.LevelDebug,
					PC:      pc,
				}
				r.Add("foo", "bar", "baz", "qux")
				return r
			},
			wantCap: 64 + len("custom log message") + 2*32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.makeRecord()
			got := NewBuilder(tt.opts, r)

			if got.sb.Cap() < tt.wantCap {
				t.Errorf("Grow() cap mismatch: got %d, want >= %d", got.sb.Cap(), tt.wantCap)
			}

			wantArgCount := r.NumAttrs()
			if tt.opts.CallerEnabled[r.Level] {
				wantArgCount++
			}

			want := &Builder{
				r:        r,
				spacePad: got.spacePad,
				argCount: wantArgCount,
				style:    got.style,
				opts:     tt.opts,
			}

			if diff := cmp.Diff(
				want,
				got,
				cmpopts.IgnoreFields(Builder{}, "sb"),
				cmp.AllowUnexported(Builder{}, slog.Record{}, pterm.RGBStyle{}),
			); diff != "" {
				t.Errorf("NewBuilder() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func removeANSI(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
}
