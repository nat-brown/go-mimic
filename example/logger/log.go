package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// New creates a new logger instance.
var New = func() logrus.FieldLogger {
	return newInnerLogger()
}

func newInnerLogger() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel
	return log
}
