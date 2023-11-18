package logger

import (
	"io"
	"log/slog"
)

const (
	envDevelopment = "development"
	envLocal       = "local"
	envProduction  = "production"
)

// New retruns a slog.Logger with text or json handler.
//
// Env can be in 3 states:
//   - development - json format, level debug
//   - local - text format, level debug
//   - production - json format, level info
//
// If env not one of this states return with text format and level info.
func New(env string, out io.Writer) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envDevelopment:
		log = slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envLocal:
		log = slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProduction:
		log = slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
