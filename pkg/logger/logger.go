package logger

import (
	"context"
	"sync"

	reqContext "github.com/eugeneradionov/froxy/pkg/context"
	"github.com/eugeneradionov/froxy/pkg/http/errors"
	"github.com/eugeneradionov/xerrors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *lgr
	once   = &sync.Once{}
)

const (
	requestIDKey = "request_id"
)

type Logger interface {
	Info(msg string, fields ...zapcore.Field)
	Debug(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)

	LogXError(ctx context.Context, xerr xerrors.XError, message string, fields ...zap.Field)
}

type lgr struct {
	logger *zap.Logger
}

func (l *lgr) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

func (l *lgr) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *lgr) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

func (l *lgr) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}

// Load initializes logger on startup.
func Load(logPreset string) (err error) {
	once.Do(func() {
		cfg := newConfig(loadConfig(logPreset))

		var zapLogger *zap.Logger

		zapLogger, err = cfg.Build()
		if err != nil {
			return
		}

		logger = &lgr{
			logger: zapLogger,
		}
	})

	return err
}

// Get returns initialized logger.
func Get() Logger {
	return logger
}

// LogXError logs XError according to it's HTTP status code.
// 0 - 499 - Info log level;
// 500+ - Error log level.
func (l *lgr) LogXError(ctx context.Context, xErr xerrors.XError, message string, fields ...zap.Field) {
	fields = append(
		fields, zap.Error(xErr),
		zap.Any("extra", xErr.GetExtra()),
		zap.Any("internal_extra", xErr.GetExtra()),
	)

	if errors.GetHTTPCode(xErr) >= 500 {
		l.withCtxValue(ctx).Error(message, fields...)
		return
	}

	l.withCtxValue(ctx).Info(message, fields...)
}

func (l *lgr) withCtxValue(ctx context.Context) *zap.Logger {
	return logger.logger.With(logger.keyAndValueFromContext(ctx)...)
}

func (l *lgr) keyAndValueFromContext(ctx context.Context) []zap.Field {
	return []zap.Field{
		zap.String(requestIDKey, reqContext.GetRequestID(ctx)),
	}
}
