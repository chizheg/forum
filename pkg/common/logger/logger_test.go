package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		expectedLevel zapcore.Level
	}{
		{
			name:          "debug level",
			level:         "debug",
			expectedLevel: zap.DebugLevel,
		},
		{
			name:          "info level",
			level:         "info",
			expectedLevel: zap.InfoLevel,
		},
		{
			name:          "warn level",
			level:         "warn",
			expectedLevel: zap.WarnLevel,
		},
		{
			name:          "error level",
			level:         "error",
			expectedLevel: zap.ErrorLevel,
		},
		{
			name:          "default level",
			level:         "invalid",
			expectedLevel: zap.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.level)
			assert.NoError(t, err)
			assert.NotNil(t, logger)
		})
	}
}

func TestLoggerMethods(t *testing.T) {
	// Создаем тестовый логгер
	observedLogs := zaptest.NewObserver()
	testLogger := &Logger{
		Logger: zap.New(observedLogs),
	}

	tests := []struct {
		name    string
		logFunc func(msg string, fields ...zap.Field)
		level   zapcore.Level
		message string
		fields  []zap.Field
	}{
		{
			name:    "debug message",
			logFunc: testLogger.Debug,
			level:   zap.DebugLevel,
			message: "debug test",
			fields:  []zap.Field{zap.String("key", "value")},
		},
		{
			name:    "info message",
			logFunc: testLogger.Info,
			level:   zap.InfoLevel,
			message: "info test",
			fields:  []zap.Field{zap.Int("count", 42)},
		},
		{
			name:    "warn message",
			logFunc: testLogger.Warn,
			level:   zap.WarnLevel,
			message: "warn test",
			fields:  []zap.Field{zap.Bool("active", true)},
		},
		{
			name:    "error message",
			logFunc: testLogger.Error,
			level:   zap.ErrorLevel,
			message: "error test",
			fields:  []zap.Field{zap.Error(assert.AnError)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observedLogs.Reset()
			tt.logFunc(tt.message, tt.fields...)

			logs := observedLogs.All()
			assert.Len(t, logs, 1)
			assert.Equal(t, tt.level, logs[0].Level)
			assert.Equal(t, tt.message, logs[0].Message)

			for _, field := range tt.fields {
				assert.Contains(t, logs[0].Context, field)
			}
		})
	}
}
