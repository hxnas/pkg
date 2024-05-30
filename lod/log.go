package lod

import "log/slog"

func ErrDebug(err error) slog.Level {
	return Iif(err == nil, slog.LevelDebug, slog.LevelWarn)
}

func ErrInfo(err error) slog.Level {
	return Iif(err == nil, slog.LevelInfo, slog.LevelWarn)
}
