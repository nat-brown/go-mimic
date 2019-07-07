package test

import (
	"strings"
	"testing"

	buffaloLogger "github.com/gobuffalo/logger"
	"github.com/sirupsen/logrus"
)

type hook struct {
	t    *testing.T
	logs []string
}

// Levels cares about logs at every level for tests.
func (h *hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire appends the log to our list so we can print at
// the end of test and avoid printing the address of
// the log file with every call.
func (h *hook) Fire(e *logrus.Entry) error {
	s, err := e.String()
	if err != nil {
		return err
	}
	h.logs = append(h.logs, s)
	return nil
}

// Print should be called outside the logrus instance
// to print the accumulated logs.
func (h *hook) Print() {
	h.t.Log("\n" + strings.Join(h.logs, ""))
}

// BlackHoleWriter prevents logrus from printing outside
// the hook.
type blackHoleWriter struct{}

// Write writes nothing, but looks like it does.
func (blackHoleWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

var logMap = map[string]*logrus.Logger{}

type logger struct {
	buffaloLogger.FieldLogger
}

func setupBuffaloLogger(t *testing.T) (buffaloLogger.FieldLogger, *hook) {
	l, h := setupLogger(t)
	return buffaloLogger.Logrus{
		FieldLogger: l,
	}, h
}

func setupLogger(t *testing.T) (*logrus.Logger, *hook) {
	l := logrus.New()
	l.Out = blackHoleWriter{}
	h := &hook{
		t:    t,
		logs: []string{},
	}
	l.Hooks.Add(h)
	logMap[t.Name()] = l
	return l, h
}
