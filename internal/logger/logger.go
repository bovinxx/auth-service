package logger

import (
	"log"
	"os"

	"github.com/natefinch/lumberjack"
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

func InitForTests() {
	globalLogger = zap.New(zap.NewNop().Core())
}

func InitDefaulLogger() {
	Init(getCore(getAtomicLevel()))
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

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set("info"); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}
