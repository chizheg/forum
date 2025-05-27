package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap logger
type Logger struct {
	*zap.Logger
}

// New creates a new logger instance
func New(debug bool) (*Logger, error) {
	var config zap.Config
	if debug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{
		Logger: logger,
	}, nil
}

// With creates a child logger with the given fields
func (l *Logger) With(fields ...zapcore.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
	}
} 