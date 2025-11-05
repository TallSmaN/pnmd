package pnmd

import (
	"log/slog"
	"reflect"
	"testing"
)

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
			cfg = Options{}
			global = nil

			Configure(tt.args.o)

			if cfg.Level != tt.wantLevel {
				t.Fatalf("Configure() cfg.Level = %v, want %v", cfg.Level, tt.wantLevel)
			}
			if !reflect.DeepEqual(cfg.CallerEnabled, tt.wantMap) {
				t.Fatalf("Configure() cfg.CallerEnabled = %#v, want %#v", cfg.CallerEnabled, tt.wantMap)
			}
			if global == nil {
				t.Fatalf("Configure() did not initialize global logger")
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
			cfg = Options{
				Level:         slog.LevelInfo,
				CallerEnabled: map[slog.Level]bool{},
			}
			for k, v := range tt.before {
				cfg.CallerEnabled[k] = v
			}
			global = nil

			DisableCallerFor(tt.args.levels...)

			if !reflect.DeepEqual(cfg.CallerEnabled, tt.after) {
				t.Fatalf("DisableCallerFor() cfg.CallerEnabled = %#v, want %#v", cfg.CallerEnabled, tt.after)
			}
			if global == nil {
				t.Fatalf("DisableCallerFor() did not rebuild global logger")
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
			cfg = Options{
				Level:         slog.LevelInfo,
				CallerEnabled: map[slog.Level]bool{},
			}
			for k, v := range tt.before {
				cfg.CallerEnabled[k] = v
			}
			global = nil

			EnableCallerFor(tt.args.levels...)

			if !reflect.DeepEqual(cfg.CallerEnabled, tt.after) {
				t.Fatalf("EnableCallerFor() cfg.CallerEnabled = %#v, want %#v", cfg.CallerEnabled, tt.after)
			}
			if global == nil {
				t.Fatalf("EnableCallerFor() did not rebuild global logger")
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		want *slog.Logger
	}{
		{
			name: "initializes global once and returns same instance",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg = Options{}
			global = nil

			tt.want = Get()
			if tt.want == nil {
				t.Fatalf("Get() returned nil on first call")
			}

			if got := Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() second call returned different instance: %p vs %p", got, tt.want)
			}
		})
	}
}

func TestSetLevel(t *testing.T) {
	type args struct {
		level slog.Level
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "set level to debug",
			args: args{level: slog.LevelDebug},
		},
		{
			name: "set level to warn",
			args: args{level: slog.LevelWarn},
		},
		{
			name: "set level to error",
			args: args{level: slog.LevelError},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg = Options{
				Level: slog.LevelInfo,
				CallerEnabled: map[slog.Level]bool{
					slog.LevelDebug: true,
					slog.LevelInfo:  true,
					slog.LevelWarn:  true,
					slog.LevelError: true,
				},
			}
			global = nil

			SetLevel(tt.args.level)

			if cfg.Level != tt.args.level {
				t.Fatalf("SetLevel() cfg.Level = %v, want %v", cfg.Level, tt.args.level)
			}
			if global == nil {
				t.Fatalf("SetLevel() did not rebuild global logger")
			}
		})
	}
}
