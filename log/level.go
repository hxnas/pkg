package log

import (
	"log/slog"
)

func Level(logger *slog.Logger, level string) {
	if h, ok := logger.Handler().(*handler); ok {
		switch level {
		case "debug":
			h.level = slog.LevelDebug
		case "info", "information":
			h.level = slog.LevelInfo
		case "warn", "warning":
			h.level = slog.LevelWarn
		case "error", "err":
			h.level = slog.LevelError
		}
	}
}

func LevelFromString(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error", "err":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func ForDefault(level string) {
	l := LevelFromString(level)
	slog.SetDefault(New(&Options{Level: l, AddSource: l == slog.LevelDebug}))
}
