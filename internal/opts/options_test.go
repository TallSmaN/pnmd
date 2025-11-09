package opts

import (
	"log/slog"
	"reflect"
	"testing"
)

func TestDefault(t *testing.T) {
	tests := []struct {
		name string
		want *Options
	}{
		{
			name: "options is default",
			want: &Options{
				Level: slog.LevelInfo,
				CallerEnabled: map[slog.Level]bool{
					slog.LevelDebug: true,
					slog.LevelInfo:  true,
					slog.LevelWarn:  true,
					slog.LevelError: true,
				},
				TimeFormat: "2006/01/02 15:04:05",
				Padding:    3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Default(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Default() = %v, want %v", got, tt.want)
			}
		})
	}
}
