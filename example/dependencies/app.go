package dependencies

import (
	"sync"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/x/sessions"

	"github.com/nat-brown/go-mimic/example/dependencies/auth"
	"github.com/nat-brown/go-mimic/example/dependencies/class"
	"github.com/nat-brown/go-mimic/example/dependencies/middleware"
	"github.com/nat-brown/go-mimic/example/dependencies/store"
)

// ENV is the app environment.
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var once = sync.Once{}

// New creates a new app.
func New() *buffalo.App {
	once.Do(func() {
		app = Construct(Options())
	})

	return app
}

// Options are the default app options.
// Returned via a function to prevent the original from being editable.
func Options() buffalo.Options {
	return buffalo.Options{
		Env:          ENV,
		SessionStore: sessions.Null{},
		SessionName:  "_deps_session",
	}
}

// Construct creates a new app but does not set it.
// Use 'New' under most (i.e. non-testing) circumstances.
func Construct(options buffalo.Options) *buffalo.App {
	app := buffalo.New(options)

	// Log request parameters (filters apply).
	app.Use(paramlogger.ParameterLogger)

	// Set the request content type to JSON.
	app.Use(contenttype.Set("application/json"))

	app.Use(middleware.ForTesting)

	app.GET("/class/recipes", class.HandleRecipes)
	app.GET("/store/ingredients/{ingredient}", store.Handle)
	app.POST("/auth/token", auth.HandleAuthentication)
	app.POST("/auth/verify", auth.HandleVerification)

	return app
}
