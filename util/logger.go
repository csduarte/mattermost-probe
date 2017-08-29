package util

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
)

// NewFileLogger will create a zap loger with essential config changes
func NewFileLogger(path string) (*logrus.Logger, error) {
	if err := ensureFile(path); err != nil {
		msg := fmt.Sprintf("Failed to find/create log file at %q", path)
		return nil, errors.Wrap(err, msg)
	}
	log := logrus.New()
	log.Level = logrus.InfoLevel

	fileHook := lfshook.NewHook(lfshook.PathMap{
		logrus.ErrorLevel: path,
		logrus.InfoLevel:  path,
		logrus.DebugLevel: path,
		logrus.PanicLevel: path,
		logrus.FatalLevel: path,
		logrus.WarnLevel:  path,
	})
	fileHook.SetFormatter(&logrus.JSONFormatter{})
	log.Hooks.Add(fileHook)

	return log, nil
}

func ensureFile(path string) error {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
