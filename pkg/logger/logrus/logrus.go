package logrus

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/demoManito/pulse/pkg/logger"
)

// logrusLogger wraps a logrus.Entry to implement logger.Logger.
type logrusLogger struct {
	entry *logrus.Entry
}

// New creates a new logrus-backed Logger.
func New(opts ...Option) logger.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})
	l.SetLevel(logrus.InfoLevel)

	for _, opt := range opts {
		opt(l)
	}

	return &logrusLogger{entry: logrus.NewEntry(l)}
}

// Option configures the underlying logrus.Logger.
type Option func(*logrus.Logger)

// WithLevel sets the log level.
func WithLevel(level logger.Level) Option {
	return func(l *logrus.Logger) {
		l.SetLevel(toLogrusLevel(level))
	}
}

// WithFormatter sets a custom logrus formatter.
func WithFormatter(f logrus.Formatter) Option {
	return func(l *logrus.Logger) {
		l.SetFormatter(f)
	}
}

func toLogrusLevel(level logger.Level) logrus.Level {
	switch level {
	case logger.DebugLevel:
		return logrus.DebugLevel
	case logger.InfoLevel:
		return logrus.InfoLevel
	case logger.WarnLevel:
		return logrus.WarnLevel
	case logger.ErrorLevel:
		return logrus.ErrorLevel
	case logger.FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func (l *logrusLogger) Debug(args ...any)           { l.entry.Debug(args...) }
func (l *logrusLogger) Debugf(f string, args ...any) { l.entry.Debugf(f, args...) }
func (l *logrusLogger) Info(args ...any)             { l.entry.Info(args...) }
func (l *logrusLogger) Infof(f string, args ...any)  { l.entry.Infof(f, args...) }
func (l *logrusLogger) Warn(args ...any)             { l.entry.Warn(args...) }
func (l *logrusLogger) Warnf(f string, args ...any)  { l.entry.Warnf(f, args...) }
func (l *logrusLogger) Error(args ...any)            { l.entry.Error(args...) }
func (l *logrusLogger) Errorf(f string, args ...any) { l.entry.Errorf(f, args...) }
func (l *logrusLogger) Fatal(args ...any)            { l.entry.Fatal(args...) }
func (l *logrusLogger) Fatalf(f string, args ...any) { l.entry.Fatalf(f, args...) }

func (l *logrusLogger) WithField(key string, value any) logger.Logger {
	return &logrusLogger{entry: l.entry.WithField(key, value)}
}

func (l *logrusLogger) WithFields(fields logger.Fields) logger.Logger {
	return &logrusLogger{entry: l.entry.WithFields(logrus.Fields(fields))}
}

func (l *logrusLogger) WithContext(ctx context.Context) logger.Logger {
	return &logrusLogger{entry: l.entry.WithContext(ctx)}
}
