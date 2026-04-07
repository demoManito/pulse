package logger

import "context"

// Fields type, used to pass to `WithFields`.
type Fields map[string]any

// Logger is the unified logging interface.
type Logger interface {
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)

	WithField(key string, value any) Logger
	WithFields(fields Fields) Logger
	WithContext(ctx context.Context) Logger
}

// defaultLogger is the package-level logger instance.
var defaultLogger Logger = nopLogger{}

// SetDefault sets the global default logger.
func SetDefault(l Logger) {
	defaultLogger = l
}

// Default returns the current default logger.
func Default() Logger {
	return defaultLogger
}

// --- package-level convenience functions ---

func Debug(args ...any)                      { defaultLogger.Debug(args...) }
func Debugf(format string, args ...any)      { defaultLogger.Debugf(format, args...) }
func Info(args ...any)                       { defaultLogger.Info(args...) }
func Infof(format string, args ...any)       { defaultLogger.Infof(format, args...) }
func Warn(args ...any)                       { defaultLogger.Warn(args...) }
func Warnf(format string, args ...any)       { defaultLogger.Warnf(format, args...) }
func Error(args ...any)                      { defaultLogger.Error(args...) }
func Errorf(format string, args ...any)      { defaultLogger.Errorf(format, args...) }
func Fatal(args ...any)                      { defaultLogger.Fatal(args...) }
func Fatalf(format string, args ...any)      { defaultLogger.Fatalf(format, args...) }
func WithField(key string, value any) Logger { return defaultLogger.WithField(key, value) }
func WithFields(fields Fields) Logger        { return defaultLogger.WithFields(fields) }
func WithContext(ctx context.Context) Logger { return defaultLogger.WithContext(ctx) }

// nopLogger is a no-op logger used as the default before initialization.
type nopLogger struct{}

func (nopLogger) Debug(...any)                         {}
func (nopLogger) Debugf(string, ...any)                {}
func (nopLogger) Info(...any)                          {}
func (nopLogger) Infof(string, ...any)                 {}
func (nopLogger) Warn(...any)                          {}
func (nopLogger) Warnf(string, ...any)                 {}
func (nopLogger) Error(...any)                         {}
func (nopLogger) Errorf(string, ...any)                {}
func (nopLogger) Fatal(...any)                         {}
func (nopLogger) Fatalf(string, ...any)                {}
func (n nopLogger) WithField(string, any) Logger       { return n }
func (n nopLogger) WithFields(Fields) Logger           { return n }
func (n nopLogger) WithContext(context.Context) Logger { return n }
