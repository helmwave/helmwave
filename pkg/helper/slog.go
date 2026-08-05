package helper

import (
	"context"
	"log/slog"
	"strings"

	log "github.com/sirupsen/logrus"
)

// SlogHandler is a slog.Handler that forwards records to logrus.
//
// helm logs through log/slog. This handler keeps helm's output inside helmwave's single log
// stream, so `--log-level` and `--log-format` control it.
type SlogHandler struct {
	// attrs are the attributes accumulated by WithAttrs, already qualified with their group prefix.
	attrs []slog.Attr
	// groups is the group prefix currently in effect, set by WithGroup.
	groups []string
}

// NewSlogHandler returns a slog.Handler writing to the standard logrus logger.
func NewSlogHandler() *SlogHandler {
	return &SlogHandler{}
}

func init() {
	// Route helm v4's package-level slog output into logrus too. NewCfg installs this handler on
	// the action.Configuration, but helm's kube waiters (pkg/kube wait.go, ready.go) log through
	// the global slog.Default() instead -- including slog.Error for failed pods and jobs. Without
	// setting the default those records skip helmwave's log stream entirely and land on Go's own
	// stderr handler, which also filters DEBUG so --progress could never surface readiness logs.
	slog.SetDefault(slog.New(NewSlogHandler()))
}

// slogLevel converts a slog level to its logrus counterpart.
func slogLevel(level slog.Level) log.Level {
	switch {
	case level < slog.LevelInfo:
		return log.DebugLevel
	case level < slog.LevelWarn:
		return log.InfoLevel
	case level < slog.LevelError:
		return log.WarnLevel
	default:
		return log.ErrorLevel
	}
}

// effectiveLevel maps a slog level to its logrus counterpart, promoting DEBUG to INFO while
// Helm.Debug is set: --progress sets Helm.Debug and is documented to surface helm's progress
// output, which helm logs through slog at DEBUG. The same flag drives downloader dependency
// progress via EnvSettings.Debug.
func effectiveLevel(level slog.Level) log.Level {
	lvl := slogLevel(level)
	if Helm.Debug && lvl == log.DebugLevel {
		return log.InfoLevel
	}

	return lvl
}

// Enabled implements slog.Handler.
func (h *SlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return log.IsLevelEnabled(effectiveLevel(level))
}

// Handle implements slog.Handler.
//
// slog.Handler says so; neither signature is ours to change.
//
//nolint:gocritic // slog.Record is passed by value and slog.Attr ranged over by value because
func (h *SlogHandler) Handle(_ context.Context, r slog.Record) error {
	fields := make(log.Fields, len(h.attrs)+r.NumAttrs())

	for _, a := range h.attrs {
		flattenAttr(fields, "", a)
	}

	prefix := strings.Join(h.groups, ".")
	r.Attrs(func(a slog.Attr) bool {
		flattenAttr(fields, prefix, a)

		return true
	})

	entry := log.WithFields(fields)
	if err, ok := fields[log.ErrorKey].(error); ok {
		entry = entry.WithError(err)
	}

	entry.Log(effectiveLevel(r.Level), r.Message)

	return nil
}

// WithAttrs implements slog.Handler.
// cheaper than the map write that follows it.
//
//nolint:gocritic // []slog.Attr is what slog.Handler hands us; copying 40 bytes per attr is
func (h *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	prefix := strings.Join(h.groups, ".")
	out := &SlogHandler{
		attrs:  make([]slog.Attr, 0, len(h.attrs)+len(attrs)),
		groups: h.groups,
	}
	out.attrs = append(out.attrs, h.attrs...)

	for _, a := range attrs {
		out.attrs = append(out.attrs, slog.Attr{Key: joinKey(prefix, a.Key), Value: a.Value})
	}

	return out
}

// WithGroup implements slog.Handler.
func (h *SlogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	out := &SlogHandler{
		attrs:  h.attrs,
		groups: make([]string, 0, len(h.groups)+1),
	}
	out.groups = append(out.groups, h.groups...)
	out.groups = append(out.groups, name)

	return out
}

// flattenAttr writes a single attribute into fields, expanding groups into dotted keys.
// slices this walks.
//
//nolint:gocritic // slog.Attr is a value type throughout slog's API, including the group
func flattenAttr(fields log.Fields, prefix string, a slog.Attr) {
	value := a.Value.Resolve()

	if value.Kind() == slog.KindGroup {
		group := value.Group()
		if len(group) == 0 {
			return
		}

		inner := prefix
		if a.Key != "" {
			inner = joinKey(prefix, a.Key)
		}

		for _, sub := range group {
			flattenAttr(fields, inner, sub)
		}

		return
	}

	if a.Key == "" {
		return
	}

	fields[joinKey(prefix, a.Key)] = value.Any()
}

// joinKey qualifies key with a group prefix.
func joinKey(prefix, key string) string {
	if prefix == "" {
		return key
	}

	return prefix + "." + key
}
