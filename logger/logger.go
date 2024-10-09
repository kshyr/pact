package logger

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/adrg/xdg"
	"github.com/sirupsen/logrus"
)

const (
	logFileRelPath = "pact/logs/demon.log"
)

type Logger struct {
	*logrus.Logger
	logFile          string
	logrusErrorFunc  func(args ...interface{})
	logrusErrorfFunc func(format string, args ...interface{})
}

func New(logFile string) (*Logger, error) {
	logger := logrus.New()
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	logger.Out = io.MultiWriter(file, os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.Error()

	return &Logger{logger, logFile, logger.Error, logger.Errorf}, nil
}

func (l *Logger) Error(args ...interface{}) error {
	l.logrusErrorFunc(args)
	return errors.New(fmt.Sprint(args))
}

func (l *Logger) Errorf(format string, args ...interface{}) error {
	l.logrusErrorfFunc(format, args...)
	return errors.New(fmt.Sprintf(format, args...))
}

func CreateLogFile(name string) (string, error) {
	relPath := fmt.Sprintf("pact/logs/%s.log", name)
	filePath, err := xdg.CacheFile(relPath)
	if err != nil {
		return "", fmt.Errorf("failed to create log file: %w", err)
	}
	os.Create(filePath)

	return filePath, nil
}
