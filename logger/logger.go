package logger

import (
	"context"

	"go.uber.org/zap"
)

type Logger struct {
	zap *zap.Logger
}

func New(level string) (*Logger, error) {
	var config zap.Config

	if level == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{zap: zapLogger}, nil
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{zap: l.zap.With(fields...)}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	requestID := GetRequestID(ctx)
	if requestID != "" {
		return l.With(zap.String("request_id", requestID))
	}
	return l
}

func (l *Logger) Sync() error {
	return l.zap.Sync()
}
