package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Encoding = "json"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	log, err := config.Build()
	if err != nil {
		panic(err)
	}

	return log
}
