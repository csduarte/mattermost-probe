package util

import (
	"github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
)

// NewFileLogger will create a zap loger with essential config changes
func NewFileLogger(path string, verbose bool) *logrus.Logger {
	log := logrus.New()
	if verbose {
		log.Level = logrus.DebugLevel
	} else {
		log.Level = logrus.InfoLevel
	}
	fileHook := lfshook.NewHook(lfshook.PathMap{
		logrus.InfoLevel: path,
	})
	fileHook.SetFormatter(&logrus.JSONFormatter{})
	log.Hooks.Add(fileHook)
	return log
}
