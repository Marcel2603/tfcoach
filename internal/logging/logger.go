package logging

import (
	"io"
	"log/slog"
	"os"
)

func SetupLogger(rawLevel string) {
	slog.SetDefault(slog.New(handler(rawLevel, os.Stderr)))
}

func handler(rawLevel string, w io.Writer) slog.Handler {
	var slogLevel slog.Level

	switch rawLevel {
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "WARN":
		slogLevel = slog.LevelWarn
	case "INFO":
		slogLevel = slog.LevelInfo
	default:
		slogLevel = slog.LevelError
	}
	logAsJSON := rawLevel == "JSON"

	if logAsJSON {
		slogLevel = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{Level: slogLevel, AddSource: true}

	if logAsJSON {
		return slog.NewJSONHandler(w, opts)
	}

	return slog.NewTextHandler(w, opts)
}
