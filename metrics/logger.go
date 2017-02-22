package metrics

import (
	"time"

	"go.uber.org/zap"
)

type LogEntry struct {
	Name            string
	time            time.Time
	DurationSeconds string
}

func NewMetricLogger() *zap.Logger {
	// // Metrics Logs
	// logConfig := zap.NewProductionConfig()
	// logConfig.DisableCaller = true
	// logConfig.Level = level
	// logConfig.OutputPaths = append(logConfig.OutputPaths, *logLocation)
	// logger, _ := logConfig.Build()
	// log = logger.Sugar()
	return nil
}
