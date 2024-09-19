package loggerinit

import (
	"log"
	"log/slog"
	"os"
)

type LogLevel string

const (
	LogInfo  LogLevel = "info"
	LogDebug LogLevel = "debug"
	LogError LogLevel = "error"
)

// setting up default slog
func MustInitLogger(level LogLevel) {
	var clog *slog.Logger
	switch level {
	case LogDebug:
		clog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case LogInfo:
		clog = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case LogError:
		clog = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	default:
		log.Fatalf("Unknown log level: %s", level)
	}
	slog.SetDefault(clog)
	slog.Info("Logger Initialized", "level", level)

}
