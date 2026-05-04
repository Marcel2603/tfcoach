package logging

import (
	"io"
	"log/slog"
	"testing"
)

func TestBuildHandler(t *testing.T) {
	tests := []struct {
		rawLevel  string
		wantLevel slog.Level
		wantJSON  bool
	}{
		{"DEBUG", slog.LevelDebug, false},
		{"INFO", slog.LevelInfo, false},
		{"WARN", slog.LevelWarn, false},
		{"ERROR", slog.LevelError, false},
		{"JSON", slog.LevelDebug, true},
		{"", slog.LevelError, false},
		{"invalid", slog.LevelError, false},
	}

	for _, tt := range tests {
		t.Run(tt.rawLevel, func(t *testing.T) {
			logHandler := handler(tt.rawLevel, io.Discard)

			_, isJSON := logHandler.(*slog.JSONHandler)
			if isJSON != tt.wantJSON {
				t.Errorf("JSON handler = %v, want %v", isJSON, tt.wantJSON)
			}

			if !logHandler.Enabled(nil, tt.wantLevel) {
				t.Errorf("level %v should be enabled", tt.wantLevel)
			}

			logLevelAbove := tt.wantLevel - 1
			if logLevelAbove < tt.wantLevel && logHandler.Enabled(nil, logLevelAbove) {
				t.Errorf("level %v should not be enabled", logLevelAbove)
			}
		})
	}
}
