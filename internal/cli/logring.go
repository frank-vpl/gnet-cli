package cli

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
)

type LogRing struct {
	mu     sync.Mutex
	lines  []string
	size   int
	latest int
}

func NewLogRing(size int) *LogRing {
	return &LogRing{lines: make([]string, size), size: size}
}

func (r *LogRing) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	line := strings.TrimRight(string(p), "\n")
	r.lines[r.latest] = line
	r.latest = (r.latest + 1) % r.size
	return len(p), nil
}

func (r *LogRing) Lines(n int) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if n > r.size {
		n = r.size
	}
	out := make([]string, 0, n)
	start := (r.latest - n + r.size) % r.size
	for i := 0; i < n; i++ {
		idx := (start + i) % r.size
		if r.lines[idx] != "" {
			out = append(out, r.lines[idx])
		}
	}
	return out
}

type ringHandler struct {
	ring  *LogRing
	inner slog.Handler
}

func (h *ringHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *ringHandler) Handle(ctx context.Context, r slog.Record) error {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s %s", r.Time.Format("15:04:05"), levelShort(r.Level), r.Message)
	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(&b, " %s=%v", a.Key, a.Value.Any())
		return true
	})
	h.ring.Write([]byte(b.String()))
	return h.inner.Handle(ctx, r)
}

func (h *ringHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ringHandler{ring: h.ring, inner: h.inner.WithAttrs(attrs)}
}

func (h *ringHandler) WithGroup(name string) slog.Handler {
	return &ringHandler{ring: h.ring, inner: h.inner.WithGroup(name)}
}

func levelShort(l slog.Level) string {
	switch l {
	case slog.LevelDebug:
		return "DBG"
	case slog.LevelInfo:
		return "INF"
	case slog.LevelWarn:
		return "WRN"
	case slog.LevelError:
		return "ERR"
	default:
		return "???"
	}
}
