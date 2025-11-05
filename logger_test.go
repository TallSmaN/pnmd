package main

import (
	"log/slog"
	"reflect"
	"testing"
)

func newTestLogger() *Logger {
	return &Logger{
		cfg: defaultLogger.cfg,
	}
}

func TestConfigure(t *testing.T) {
	type args struct {
		o Options
	}
	tests := []struct {
		name      string
		args      args
		wantLevel slog.Level
		wantMap   map[slog.Level]bool
	}{
		{
			name: "empty options -> normalized defaults",
			args: args{
				o: Options{},
			},
			wantLevel: slog.LevelInfo,
			wantMap: map[slog.Level]bool{
				slog.LevelDebug: true,
				slog.LevelInfo:  true,
				slog.LevelWarn:  true,
				slog.LevelError: true,
			},
		},
		{
			name: "custom level warn, nil map -> default caller map",
			args: args{
				o: Options{
					Level:         slog.LevelWarn,
					CallerEnabled: nil,
				},
			},
			wantLevel: slog.LevelWarn,
			wantMap: map[slog.Level]bool{
				slog.LevelDebug: true,
				slog.LevelInfo:  true,
				slog.LevelWarn:  true,
				slog.LevelError: true,
			},
		},
		{
			name: "custom caller map overrides default",
			args: args{
				o: Options{
					Level: slog.LevelDebug,
					CallerEnabled: map[slog.Level]bool{
						slog.LevelDebug: true,
						slog.LevelInfo:  false,
						slog.LevelWarn:  true,
						slog.LevelError: false,
					},
				},
			},
			wantLevel: slog.LevelDebug,
			wantMap: map[slog.Level]bool{
				slog.LevelDebug: true,
				slog.LevelInfo:  false,
				slog.LevelWarn:  true,
				slog.LevelError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newTestLogger()

			l.Configure(tt.args.o)

			if l.cfg.Level != tt.wantLevel {
				t.Fatalf("Configure() cfg.Level = %v, want %v", l.cfg.Level, tt.wantLevel)
			}
			if !reflect.DeepEqual(l.cfg.CallerEnabled, tt.wantMap) {
				t.Fatalf("Configure() cfg.CallerEnabled = %#v, want %#v", l.cfg.CallerEnabled, tt.wantMap)
			}
			if l.Logger == nil {
				t.Fatalf("Configure() did not initialize underlying slog.Logger")
			}
		})
	}
}

func TestDisableCallerFor(t *testing.T) {
	type args struct {
		levels []slog.Level
	}
	tests := []struct {
		name   string
		args   args
		before map[slog.Level]bool
		after  map[slog.Level]bool
	}{
		{
			name: "disable debug and info",
			args: args{
				levels: []slog.Level{slog.LevelDebug, slog.LevelInfo},
			},
			before: map[slog.Level]bool{
				slog.LevelDebug: true,
				slog.LevelInfo:  true,
				slog.LevelWarn:  true,
				slog.LevelError: true,
			},
			after: map[slog.Level]bool{
				slog.LevelDebug: false,
				slog.LevelInfo:  false,
				slog.LevelWarn:  true,
				slog.LevelError: true,
			},
		},
		{
			name: "disable error only",
			args: args{
				levels: []slog.Level{slog.LevelError},
			},
			before: map[slog.Level]bool{
				slog.LevelDebug: true,
				slog.LevelInfo:  true,
				slog.LevelWarn:  true,
				slog.LevelError: true,
			},
			after: map[slog.Level]bool{
				slog.LevelDebug: true,
				slog.LevelInfo:  true,
				slog.LevelWarn:  true,
				slog.LevelError: false,
			},
		},
		{
			name: "idempotent when called twice",
			args: args{
				levels: []slog.Level{slog.LevelWarn, slog.LevelWarn},
			},
			before: map[slog.Level]bool{
				slog.LevelDebug: true,
				slog.LevelInfo:  true,
				slog.LevelWarn:  true,
				slog.LevelError: true,
			},
			after: map[slog.Level]bool{
				slog.LevelDebug: true,
				slog.LevelInfo:  true,
				slog.LevelWarn:  false,
				slog.LevelError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newTestLogger()
			l.cfg.Level = slog.LevelInfo
			l.cfg.CallerEnabled = map[slog.Level]bool{}

			for k, v := range tt.before {
				l.cfg.CallerEnabled[k] = v
			}

			l.Logger = nil

			l.DisableCallerFor(tt.args.levels...)

			if !reflect.DeepEqual(l.cfg.CallerEnabled, tt.after) {
				t.Fatalf("DisableCallerFor() cfg.CallerEnabled = %#v, want %#v", l.cfg.CallerEnabled, tt.after)
			}
			if l.Logger == nil {
				t.Fatalf("DisableCallerFor() did not rebuild underlying slog.Logger")
			}
		})
	}
}

func TestEnableCallerFor(t *testing.T) {
	type args struct {
		levels []slog.Level
	}
	tests := []struct {
		name   string
		args   args
		before map[slog.Level]bool
		after  map[slog.Level]bool
	}{
		{
			name: "enable info and error",
			args: args{
				levels: []slog.Level{slog.LevelInfo, slog.LevelError},
			},
			before: map[slog.Level]bool{
				slog.LevelDebug: false,
				slog.LevelInfo:  false,
				slog.LevelWarn:  false,
				slog.LevelError: false,
			},
			after: map[slog.Level]bool{
				slog.LevelDebug: false,
				slog.LevelInfo:  true,
				slog.LevelWarn:  false,
				slog.LevelError: true,
			},
		},
		{
			name: "idempotent when already true",
			args: args{
				levels: []slog.Level{slog.LevelWarn},
			},
			before: map[slog.Level]bool{
				slog.LevelDebug: false,
				slog.LevelInfo:  true,
				slog.LevelWarn:  true,
				slog.LevelError: false,
			},
			after: map[slog.Level]bool{
				slog.LevelDebug: false,
				slog.LevelInfo:  true,
				slog.LevelWarn:  true,
				slog.LevelError: false,
			},
		},
		{
			name: "enable debug only",
			args: args{
				levels: []slog.Level{slog.LevelDebug},
			},
			before: map[slog.Level]bool{
				slog.LevelDebug: false,
				slog.LevelInfo:  true,
				slog.LevelWarn:  false,
				slog.LevelError: true,
			},
			after: map[slog.Level]bool{
				slog.LevelDebug: true,
				slog.LevelInfo:  true,
				slog.LevelWarn:  false,
				slog.LevelError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newTestLogger()
			l.cfg.Level = slog.LevelInfo
			l.cfg.CallerEnabled = map[slog.Level]bool{}

			for k, v := range tt.before {
				l.cfg.CallerEnabled[k] = v
			}

			l.Logger = nil

			l.EnableCallerFor(tt.args.levels...)

			if !reflect.DeepEqual(l.cfg.CallerEnabled, tt.after) {
				t.Fatalf("EnableCallerFor() cfg.CallerEnabled = %#v, want %#v", l.cfg.CallerEnabled, tt.after)
			}
			if l.Logger == nil {
				t.Fatalf("EnableCallerFor() did not rebuild underlying slog.Logger")
			}
		})
	}
}

func TestGet_PackageLoggerInstance(t *testing.T) {
	defaultLogger = Get()

	l1 := Get()
	if l1 == nil || l1.Logger == nil {
		t.Fatalf("Get() must build underlying *slog.Logger on first call")
	}
	core1 := l1.Logger

	l2 := Get()
	if l2 != l1 {
		t.Errorf("Get() must return the same *Logger instance: %p vs %p", l1, l2)
	}
	if l2.Logger != core1 {
		t.Errorf("core must be the same before reconfigure")
	}

	l2.Configure(Options{})
	core2 := Get().Logger
	if core2 == nil {
		t.Fatalf("nil core after Configure")
	}
	if core2 == core1 {
		t.Errorf("Configure() must rebuild underlying *slog.Logger")
	}

	beforeCore := Get().Logger
	defaultLogger.SetLevel(slog.LevelError)
	afterCore := Get().Logger
	if afterCore == nil {
		t.Fatalf("nil core after SetLevel")
	}
	if afterCore == beforeCore {
		t.Errorf("SetLevel() must rebuild underlying *slog.Logger")
	}
}

func TestSetLevel(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
	}{
		{name: "set level to debug", level: slog.LevelDebug},
		{name: "set level to warn", level: slog.LevelWarn},
		{name: "set level to error", level: slog.LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newTestLogger()
			prev := l.get()

			l.SetLevel(tt.level)

			if l.cfg.Level != tt.level {
				t.Fatalf("SetLevel() cfg.Level = %v, want %v", l.cfg.Level, tt.level)
			}
			if l.Logger == nil {
				t.Fatalf("SetLevel() did not rebuild underlying slog.Logger")
			}
			if l.Logger == prev {
				t.Errorf("SetLevel() should rebuild slog.Logger instance")
			}
		})
	}
}
