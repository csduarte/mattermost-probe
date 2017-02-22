package util

import (
	"log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = SetupLogger()

	if err != nil {
		log.Fatal(err)
	}
}

func SetupLogger() (*zap.Logger, error) {
	logConfig := zap.NewProductionConfig()
	logConfig.Encoding = "console"
	logConfig.DisableCaller = true
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := logConfig.Build()
	return logger, err
}

func LogError(message string, fields ...zapcore.Field) {
	logger.Error(message, fields...)
}

func LogInfo(message string, fields ...zapcore.Field) {
	logger.Info(message, fields...)
}