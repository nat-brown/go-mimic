package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	// HeaderContentType header key for content type.
	HeaderContentType = "Content-Type"
	// HeaderAccept header key for accept.
	HeaderAccept = "Accept"

	// AcceptJSON header value for accept json.
	AcceptJSON = "application/json"
	// AcceptText header value for accept plain text.
	AcceptText = "text/plain"
	// AcceptUnspecified header value for accept was not populated.
	AcceptUnspecified = ""
)

// New returns a new handler.
func New(client *http.Client, logger logrus.FieldLogger) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/recipes", RecipeHandler{
		Client: client,
		Logger: logger,
	}.Handle).Methods(http.MethodGet)

	return r
}
