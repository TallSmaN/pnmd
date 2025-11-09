package handler

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/TallSmaN/pnmd/internal/opts"
)

func TestHandler_Enabled(t *testing.T) {
	type fields struct {
		Opts *opts.Options
		W    io.Writer
	}

	type args struct {
		ctx   context.Context
		level slog.Level
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "enabled when level >= configured",
			fields: fields{
				Opts: &opts.Options{
					Level: slog.LevelWarn,
					CallerEnabled: map[slog.Level]bool{
						slog.LevelDebug: true,
						slog.LevelInfo:  true,
						slog.LevelWarn:  true,
						slog.LevelError: true,
					},
					TimeFormat: "2006/01/02 15:04:05",
					Padding:    3,
				},
				W: os.Stdout,
			},
			args: args{
				ctx:   context.Background(),
				level: slog.LevelWarn,
			},
			want: true,
		},
		{
			name: "disabled when level < configured",
			fields: fields{
				Opts: &opts.Options{
					Level: slog.LevelWarn,
					CallerEnabled: map[slog.Level]bool{
						slog.LevelDebug: true,
						slog.LevelInfo:  true,
						slog.LevelWarn:  true,
						slog.LevelError: true,
					},
					TimeFormat: "2006/01/02 15:04:05",
					Padding:    3,
				},
				W: os.Stdout,
			},
			args: args{
				ctx:   context.Background(),
				level: slog.LevelInfo,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				Opts: tt.fields.Opts,
				W:    tt.fields.W,
			}

			if got := h.Enabled(tt.args.ctx, tt.args.level); got != tt.want {
				t.Errorf("Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
