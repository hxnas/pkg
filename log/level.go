package log

import (
	"log/slog"
	"strings"
)

// Level sets the log level
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

func AddSource(logger *slog.Logger, debugOnly bool) {
	if h, ok := logger.Handler().(*handler); ok {
		h.addSource = true
		h.debugSourceOnly = debugOnly
	}
}

// LevelFromString parse the level from string, ignore case
//   - debug => slog.LevelDebug
//   - info, information => slog.LevelInfo
//   - warn, warning => slog.LevelWarn
//   - error, err => slog.LevelError
//   - otherwise slog.LevelInfo
func LevelFromString(level string) slog.Level {
	switch strings.ToLower(level) {
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

// ForDefault 函数用于设置默认日志记录器的级别和是否添加源信息。
//
// 参数:
//   - level string - 日志级别的字符串表示，例如 "info", "debug" 等。
//   - addSource ...bool - 可选参数，用于指定是否添加源信息。如果提供，第一个布尔值将被用于设置是否添加源信息。否则根据 level == debug 判断添加源信息。
func ForDefault(level string, addSource ...bool) {
	l := LevelFromString(level)
	// 初始化选项，当日志级别为Debug时，默认添加源信息
	lOpt := &Options{Level: l, AddSource: l == slog.LevelDebug}
	// 如果提供了addSource参数，则使用提供的值覆盖默认设置
	if len(addSource) > 0 {
		lOpt.AddSource = addSource[0]
	}
	// 设置默认日志记录器
	slog.SetDefault(New(lOpt))
}

func SetDefault(opts *Options) {
	slog.SetDefault(New(opts))
}
