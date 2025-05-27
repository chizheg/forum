package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(level string) (*Logger, error) {
	config := zap.NewProductionConfig()

	// Установка уровня логирования
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Настройка формата вывода
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logger, err := config.Build(
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{
		Logger: logger,
	}, nil
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}
