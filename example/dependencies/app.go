package dependencies

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/x/sessions"

	"github.com/nat-brown/go-mimic/example/dependencies/auth"
	"github.com/nat-brown/go-mimic/example/dependencies/class"
	"github.com/nat-brown/go-mimic/example/dependencies/store"
)

// ENV is the app environment.
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// New creates a new app.
func New() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			SessionName:  "_deps_session",
		})

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Set the request content type to JSON
		app.Use(contenttype.Set("application/json"))

		app.GET("/class/recipes", class.HandleRecipes)
		app.GET("/store/ingredients/{ingredient}", store.Handle)
		app.POST("/auth/token", auth.HandleAuthentication)
		app.POST("/auth/verify", auth.HandleVerification)
	}

	return app
}
