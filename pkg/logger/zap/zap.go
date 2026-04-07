package zap

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/demoManito/pulse/pkg/logger"
)

// zapLogger wraps a zap.SugaredLogger to implement logger.Logger.
type zapLogger struct {
	sugar  *zap.SugaredLogger
	ctx    context.Context
	fields logger.Fields
}

// New creates a new zap-backed Logger.
func New(opts ...Option) logger.Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	for _, opt := range opts {
		opt(&cfg)
	}

	l, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	return &zapLogger{sugar: l.Sugar()}
}

// Option configures the underlying zap config.
type Option func(*zap.Config)

// WithLevel sets the log level.
func WithLevel(level logger.Level) Option {
	return func(cfg *zap.Config) {
		cfg.Level = zap.NewAtomicLevelAt(toZapLevel(level))
	}
}

// WithDevelopment enables development mode (console encoding, debug level, stack traces).
func WithDevelopment() Option {
	return func(cfg *zap.Config) {
		cfg.Development = true
		cfg.Encoding = "console"
	}
}

func toZapLevel(level logger.Level) zapcore.Level {
	switch level {
	case logger.DebugLevel:
		return zapcore.DebugLevel
	case logger.InfoLevel:
		return zapcore.InfoLevel
	case logger.WarnLevel:
		return zapcore.WarnLevel
	case logger.ErrorLevel:
		return zapcore.ErrorLevel
	case logger.FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l *zapLogger) withSugar(s *zap.SugaredLogger) *zapLogger {
	return &zapLogger{sugar: s, ctx: l.ctx, fields: l.fields}
}

func (l *zapLogger) Debug(args ...any)           { l.sugar.Debug(args...) }
func (l *zapLogger) Debugf(f string, args ...any) { l.sugar.Debugf(f, args...) }
func (l *zapLogger) Info(args ...any)             { l.sugar.Info(args...) }
func (l *zapLogger) Infof(f string, args ...any)  { l.sugar.Infof(f, args...) }
func (l *zapLogger) Warn(args ...any)             { l.sugar.Warn(args...) }
func (l *zapLogger) Warnf(f string, args ...any)  { l.sugar.Warnf(f, args...) }
func (l *zapLogger) Error(args ...any)            { l.sugar.Error(args...) }
func (l *zapLogger) Errorf(f string, args ...any) { l.sugar.Errorf(f, args...) }
func (l *zapLogger) Fatal(args ...any)            { l.sugar.Fatal(args...) }
func (l *zapLogger) Fatalf(f string, args ...any) { l.sugar.Fatalf(f, args...) }

func (l *zapLogger) WithField(key string, value any) logger.Logger {
	return l.withSugar(l.sugar.With(key, value))
}

func (l *zapLogger) WithFields(fields logger.Fields) logger.Logger {
	kvs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		kvs = append(kvs, k, v)
	}
	return l.withSugar(l.sugar.With(kvs...))
}

func (l *zapLogger) WithContext(ctx context.Context) logger.Logger {
	return &zapLogger{sugar: l.sugar, ctx: ctx, fields: l.fields}
}
