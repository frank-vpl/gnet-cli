package logger

import (
	"io"
	"log/slog"
	"os"
	"sync"
)

type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

var (
	mu       sync.Mutex
	current  *slog.Logger
	fileSink io.Closer
)

func init() {
	current = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: LevelInfo}))
}

func L() *slog.Logger {
	mu.Lock()
	defer mu.Unlock()
	return current
}

func SetLevel(level Level) {
	mu.Lock()
	defer mu.Unlock()
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	current = slog.New(handler)
}

func SetHandler(handler slog.Handler) {
	mu.Lock()
	defer mu.Unlock()
	current = slog.New(handler)
}

func NewTextHandler(level Level) slog.Handler {
	return slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
}

func EnableFileLog(path string) error {
	mu.Lock()
	defer mu.Unlock()
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	fileSink = f
	mw := io.MultiWriter(os.Stderr, f)
	handler := slog.NewTextHandler(mw, &slog.HandlerOptions{Level: LevelDebug})
	current = slog.New(handler)
	return nil
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()
	if fileSink != nil {
		err := fileSink.Close()
		fileSink = nil
		return err
	}
	return nil
}

func Debug(msg string, args ...any) { L().Debug(msg, args...) }
func Info(msg string, args ...any)  { L().Info(msg, args...) }
func Warn(msg string, args ...any)  { L().Warn(msg, args...) }
func Error(msg string, args ...any) { L().Error(msg, args...) }
