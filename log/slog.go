package log

import (
	"context"
	"encoding"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
	"unicode"
)

// ANSI modes
const (
	ansiReset          = "\033[0m"
	ansiFaint          = "\033[2m"
	ansiResetFaint     = "\033[22m"
	ansiBrightRed      = "\033[91m"
	ansiBrightGreen    = "\033[92m"
	ansiBrightYellow   = "\033[93m"
	ansiBrightRedFaint = "\033[91;2m"
)

const errKey = "err"

var (
	defaultLevel      = slog.LevelInfo
	defaultTimeFormat = "01/02 15:04:05"
)

// Options for a slog.Handler that writes tinted logs. A zero Options consists
// entirely of default values.
//
// Options can be used as a drop-in replacement for [slog.HandlerOptions].
type Options struct {
	// Enable source code location (Default: false)
	AddSource bool

	// Minimum level to log (Default: slog.LevelInfo)
	Level slog.Leveler

	// ReplaceAttr is called to rewrite each non-group attribute before it is logged.
	// See https://pkg.go.dev/log/slog#HandlerOptions for details.
	ReplaceAttr func(groups []string, attr slog.Attr) slog.Attr

	// Time format (Default: "01/02 15:04:05")
	TimeFormat string

	// Disable color (Default: false)
	NoColor bool

	DebugSourceOnly bool
}

// NewHandler creates a [slog.Handler] that writes tinted logs to Writer w,
// using the default options. If opts is nil, the default options are used.
func NewHandler(w io.Writer, opts *Options) slog.Handler {
	h := &handler{
		w:          w,
		level:      defaultLevel,
		timeFormat: defaultTimeFormat,
	}
	if opts == nil {
		return h
	}

	h.addSource = opts.AddSource
	h.debugSourceOnly = opts.DebugSourceOnly

	if opts.Level != nil {
		h.level = opts.Level
	}
	h.replaceAttr = opts.ReplaceAttr
	if opts.TimeFormat != "" {
		h.timeFormat = opts.TimeFormat
	}
	h.noColor = opts.NoColor
	return h
}

func New(opts *Options) *slog.Logger {
	return slog.New(NewHandler(os.Stderr, opts))
}

// handler implements a [slog.Handler].
type handler struct {
	attrsPrefix string
	groupPrefix string
	groups      []string

	mu sync.Mutex
	w  io.Writer

	debugSourceOnly bool
	addSource       bool
	level           slog.Leveler
	replaceAttr     func([]string, slog.Attr) slog.Attr
	timeFormat      string
	noColor         bool
}

func (h *handler) clone() *handler {
	return &handler{
		attrsPrefix:     h.attrsPrefix,
		groupPrefix:     h.groupPrefix,
		groups:          h.groups,
		w:               h.w,
		addSource:       h.addSource,
		debugSourceOnly: h.debugSourceOnly,
		level:           h.level,
		replaceAttr:     h.replaceAttr,
		timeFormat:      h.timeFormat,
		noColor:         h.noColor,
	}
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

type contextKey struct{ name string }

var prefixKey = contextKey{"prefix"}

func Prefix(ctx context.Context, prefix string) context.Context {
	return context.WithValue(ctx, prefixKey, prefix)
}

func GetPrefix(ctx context.Context) (prefix string) {
	prefix, _ = ctx.Value(prefixKey).(string)
	return
}

func (h *handler) resloveRecord(r *slog.Record) {
	var sourceExist bool
	r.Attrs(func(attr slog.Attr) bool {
		if attr.Key == slog.TimeKey && attr.Value.Kind() == slog.KindTime {
			r.Time = attr.Value.Time()
		}

		switch attr.Key {
		case slog.TimeKey:
			if attr.Value.Kind() == slog.KindTime {
				if t := attr.Value.Time(); !t.IsZero() {
					r.Time = attr.Value.Time()
				}
				attr.Value = slog.AnyValue(nil)
			}
			return !h.addSource
		case slog.SourceKey:
			sourceExist = true
		}
		return true
	})

	if !sourceExist && h.addSource && (r.Level == slog.LevelDebug || !h.debugSourceOnly) {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			f.File = filepath.ToSlash(f.File)
			dir, file := path.Split(f.File)
			dir = path.Base(dir)
			filename := path.Join(path.Base(dir), file)
			src := &slog.Source{
				Function: f.Function,
				File:     filename,
				Line:     f.Line,
			}
			r.AddAttrs(slog.Any(slog.SourceKey, src))
		}
	}
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	// get a buffer from the sync pool
	buf := newBuffer()
	defer buf.Free()

	rep := h.replaceAttr

	h.resloveRecord(&r)

	// write time
	if !r.Time.IsZero() {
		val := r.Time.Round(0) // strip monotonic to match Attr behavior
		if rep == nil {
			h.appendTime(buf, r.Time)
			buf.WriteByte(' ')
		} else if a := rep(nil /* groups */, slog.Time(slog.TimeKey, val)); a.Key != "" {
			if a.Value.Kind() == slog.KindTime {
				h.appendTime(buf, a.Value.Time())
			} else {
				h.appendValue(buf, a.Value, false)
			}
			buf.WriteByte(' ')
		}
	}

	// write level
	if rep == nil {
		h.appendLevel(buf, r.Level)
		buf.WriteByte(' ')
	} else if a := rep(nil /* groups */, slog.Any(slog.LevelKey, r.Level)); a.Key != "" {
		h.appendValue(buf, a.Value, false)
		buf.WriteByte(' ')
	}

	if prefix := GetPrefix(ctx); prefix != "" {
		buf.WriteByte('[')
		buf.WriteString(prefix)
		buf.WriteByte(']')
		buf.WriteByte(' ')
	}

	// write message
	if rep == nil {
		buf.WriteString(r.Message)
		buf.WriteByte(' ')
	} else if a := rep(nil /* groups */, slog.String(slog.MessageKey, r.Message)); a.Key != "" {
		h.appendValue(buf, a.Value, false)
		buf.WriteByte(' ')
	}

	// write handler attributes
	if len(h.attrsPrefix) > 0 {
		buf.WriteString(h.attrsPrefix)
	}

	// write attributes
	r.Attrs(func(attr slog.Attr) bool {
		if attr.Key != slog.TimeKey || attr.Value.Kind() != slog.KindTime {
			if attr.Value.Any() != nil {
				h.appendAttr(buf, attr, h.groupPrefix, h.groups)
			}
		}
		return true
	})

	if len(*buf) == 0 {
		return nil
	}
	(*buf)[len(*buf)-1] = '\n' // replace last space with newline

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.w.Write(*buf)
	return err
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := h.clone()

	buf := newBuffer()
	defer buf.Free()

	// write attributes to buffer
	for _, attr := range attrs {
		h.appendAttr(buf, attr, h.groupPrefix, h.groups)
	}
	h2.attrsPrefix = h.attrsPrefix + string(*buf)
	return h2
}

func (h *handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := h.clone()
	h2.groupPrefix += name + "."
	h2.groups = append(h2.groups, name)
	return h2
}

func (h *handler) appendTime(buf *buffer, t time.Time) {
	buf.WriteStringIf(!h.noColor, ansiFaint)
	*buf = t.AppendFormat(*buf, h.timeFormat)
	buf.WriteStringIf(!h.noColor, ansiReset)
}

func (h *handler) appendLevel(buf *buffer, level slog.Level) {
	switch {
	case level < slog.LevelInfo:
		buf.WriteString("DBG")
		appendLevelDelta(buf, level-slog.LevelDebug)
	case level < slog.LevelWarn:
		buf.WriteStringIf(!h.noColor, ansiBrightGreen)
		buf.WriteString("INF")
		appendLevelDelta(buf, level-slog.LevelInfo)
		buf.WriteStringIf(!h.noColor, ansiReset)
	case level < slog.LevelError:
		buf.WriteStringIf(!h.noColor, ansiBrightYellow)
		buf.WriteString("WRN")
		appendLevelDelta(buf, level-slog.LevelWarn)
		buf.WriteStringIf(!h.noColor, ansiReset)
	default:
		buf.WriteStringIf(!h.noColor, ansiBrightRed)
		buf.WriteString("ERR")
		appendLevelDelta(buf, level-slog.LevelError)
		buf.WriteStringIf(!h.noColor, ansiReset)
	}
}

func appendLevelDelta(buf *buffer, delta slog.Level) {
	if delta == 0 {
		return
	} else if delta > 0 {
		buf.WriteByte('+')
	}
	*buf = strconv.AppendInt(*buf, int64(delta), 10)
}

func (h *handler) appendSource(buf *buffer, src *slog.Source) {
	buf.WriteStringIf(!h.noColor, ansiFaint)
	*buf = strconv.AppendQuote(*buf, src.File+":"+strconv.Itoa(src.Line))
	buf.WriteStringIf(!h.noColor, ansiReset)
}

func (h *handler) appendAttr(buf *buffer, attr slog.Attr, groupsPrefix string, groups []string) {
	attr.Value = attr.Value.Resolve()
	if rep := h.replaceAttr; rep != nil && attr.Value.Kind() != slog.KindGroup {
		attr = rep(groups, attr)
		attr.Value = attr.Value.Resolve()
	}

	if attr.Equal(slog.Attr{}) {
		return
	}

	if attr.Value.Kind() == slog.KindGroup {
		if attr.Key != "" {
			groupsPrefix += attr.Key + "."
			groups = append(groups, attr.Key)
		}
		for _, groupAttr := range attr.Value.Group() {
			h.appendAttr(buf, groupAttr, groupsPrefix, groups)
		}
	} else if err, ok := attr.Value.Any().(tintError); ok {
		// append tintError
		h.appendTintError(buf, err, groupsPrefix)
		buf.WriteByte(' ')
	} else {
		h.appendKey(buf, attr.Key, groupsPrefix)
		h.appendValue(buf, attr.Value, true)
		buf.WriteByte(' ')
	}
}

func (h *handler) appendKey(buf *buffer, key, groups string) {
	buf.WriteStringIf(!h.noColor, ansiFaint)
	appendString(buf, groups+key, true)
	buf.WriteByte('=')
	buf.WriteStringIf(!h.noColor, ansiReset)
}

func (h *handler) appendValue(buf *buffer, v slog.Value, quote bool) {
	switch v.Kind() {
	case slog.KindString:
		appendString(buf, v.String(), quote)
	case slog.KindInt64:
		*buf = strconv.AppendInt(*buf, v.Int64(), 10)
	case slog.KindUint64:
		*buf = strconv.AppendUint(*buf, v.Uint64(), 10)
	case slog.KindFloat64:
		*buf = strconv.AppendFloat(*buf, v.Float64(), 'g', -1, 64)
	case slog.KindBool:
		*buf = strconv.AppendBool(*buf, v.Bool())
	case slog.KindDuration:
		appendString(buf, v.Duration().String(), quote)
	case slog.KindTime:
		appendString(buf, v.Time().String(), quote)
	case slog.KindAny:
		switch cv := v.Any().(type) {
		case slog.Level:
			h.appendLevel(buf, cv)
		case encoding.TextMarshaler:
			data, err := cv.MarshalText()
			if err != nil {
				break
			}
			appendString(buf, string(data), quote)
		case *slog.Source:
			h.appendSource(buf, cv)
		default:
			appendString(buf, fmt.Sprintf("%+v", v.Any()), quote)
		}
	}
}

func (h *handler) appendTintError(buf *buffer, err error, groupsPrefix string) {
	buf.WriteStringIf(!h.noColor, ansiBrightRedFaint)
	appendString(buf, groupsPrefix+errKey, true)
	buf.WriteByte('=')
	buf.WriteStringIf(!h.noColor, ansiResetFaint)
	appendString(buf, err.Error(), true)
	buf.WriteStringIf(!h.noColor, ansiReset)
}

func appendString(buf *buffer, s string, quote bool) {
	if quote && needsQuoting(s) {
		*buf = strconv.AppendQuote(*buf, s)
	} else {
		buf.WriteString(s)
	}
}

func needsQuoting(s string) bool {
	if len(s) == 0 {
		return true
	}
	for _, r := range s {
		if unicode.IsSpace(r) || r == '"' || r == '=' || !unicode.IsPrint(r) {
			return true
		}
	}
	return false
}

type tintError struct{ error }

// Err returns a tinted (colorized) [slog.Attr] that will be written in red color
// by the [tint.Handler]. When used with any other [slog.Handler], it behaves as
//
//	slog.Any("err", err)
func Err(err error) slog.Attr {
	if err != nil {
		err = tintError{err}
	}
	return slog.Any(errKey, err)
}
