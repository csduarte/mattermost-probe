package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// var logger *zap.Logger

// func init() {
// 	var err error
// 	logger, err = NewLogger()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// NewFileLogger will create a zap loger with essential config changes
func NewFileLogger(path string) (*zap.Logger, error) {
	logConfig := zap.NewProductionConfig()
	level := zap.NewAtomicLevel()
	level.SetLevel(zap.DebugLevel)

	logConfig.OutputPaths = append(logConfig.OutputPaths, path)
	logConfig.Level = level
	logConfig.DisableCaller = true
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return logConfig.Build()
}

// NewEasyLogger will create a sugared logger
func NewEasyLogger(path string) (*zap.SugaredLogger, error) {
	log, err := NewFileLogger(path)
	if err != nil {
		return nil, err
	}
	return log.Sugar(), nil
}

// // LogError allows for easy global Error log
// func LogError(message string, fields ...zapcore.Field) {
// 	logger.Error(message, fields...)
// }

// // LogInfo allows for easy global Info log
// func LogInfo(message string, fields ...zapcore.Field) {
// 	logger.Info(message, fields...)
// }
