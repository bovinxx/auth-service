package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

type Field struct {
	Key   string
	Value any
}

func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...)
}

func Err(err error) Field {
	return Field{
		Key:   "error",
		Value: err,
	}
}

func convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}

	return zapFields
}

func Debug(msg string, fields ...Field) {
	globalLogger.Debug(msg, convertFields(fields)...)
}

func Info(msg string, fields ...Field) {
	globalLogger.Info(msg, convertFields(fields)...)
}

func Warn(msg string, fields ...Field) {
	globalLogger.Warn(msg, convertFields(fields)...)
}

func Error(msg string, fields ...Field) {
	globalLogger.Error(msg, convertFields(fields)...)
}

func Fatal(msg string, fields ...Field) {
	globalLogger.Fatal(msg, convertFields(fields)...)
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(opts...)
}
