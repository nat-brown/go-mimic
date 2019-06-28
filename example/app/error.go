package app

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// ErrorHandler handles errors.
type ErrorHandler interface {
	Handle(http.ResponseWriter, logrus.FieldLogger)
}

// ResponseError keeps track of errors to return to users and errors to log.
type ResponseError struct {
	InnerError   string
	ReturnError  string
	ResponseCode int
}

// Handle deals with actions associated with the given error.
func (e ResponseError) Handle(w http.ResponseWriter, l logrus.FieldLogger) {
	l.Info(e.InnerError)
	http.Error(w, e.ReturnError, e.ResponseCode)
}

// NoMessageError keeps track of errors to return to users and errors to log.
type NoMessageError struct {
	InnerError   string
	ResponseCode int
}

// Handle deals with actions associated with the given error.
func (e NoMessageError) Handle(w http.ResponseWriter, l logrus.FieldLogger) {
	l.Info(e.InnerError)
	w.WriteHeader(e.ResponseCode)
}
