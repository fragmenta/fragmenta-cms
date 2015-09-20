package app

import (
	"github.com/fragmenta/router"

	"github.com/fragmenta/fragmenta-cms/src/pages/actions"
	"github.com/fragmenta/fragmenta-cms/src/tags/actions"
	"github.com/fragmenta/fragmenta-cms/src/users/actions"
)

// Define routes for this app
func setupRoutes(r *router.Router) {

	// Set the default file handler
	r.FileHandler = fileHandler
	r.ErrorHandler = errHandler

	// Add a files route to handle static images under files - nginx deals with this in production
	r.Add("/files/{path:.*}", fileHandler)

	// Add the home page route
	r.Add("/", pageactions.HandleHome)

	// Standard REST handlers for tags
	r.Add("/tags", tagactions.HandleIndex)
	r.Add("/tags/create", tagactions.HandleCreateShow)
	r.Add("/tags/create", tagactions.HandleCreate).Post()
	r.Add("/tags/{id:[0-9]+}/update", tagactions.HandleUpdateShow)
	r.Add("/tags/{id:[0-9]+}/update", tagactions.HandleUpdate).Post()
	r.Add("/tags/{id:[0-9]+}/destroy", tagactions.HandleDestroy).Post()
	r.Add("/tags/{id:[0-9]+}", tagactions.HandleShow)

	// Standard REST handlers for users
	r.Add("/users", useractions.HandleIndex)
	r.Add("/users/create", useractions.HandleCreateShow)
	r.Add("/users/create", useractions.HandleCreate).Post()
	r.Add("/users/{id:[0-9]+}/update", useractions.HandleUpdateShow)
	r.Add("/users/{id:[0-9]+}/update", useractions.HandleUpdate).Post()
	r.Add("/users/{id:[0-9]+}/destroy", useractions.HandleDestroy).Post()
	r.Add("/users/{id:[0-9]+}", useractions.HandleShow)
	r.Add("/users/login", useractions.HandleLoginShow)
	r.Add("/users/login", useractions.HandleLogin).Post()
	r.Add("/users/logout", useractions.HandleLogout).Post()
	r.Add("/users/password", useractions.HandlePasswordReset).Post()
	r.Add("/users/password/reset", useractions.HandlePasswordResetShow)
	r.Add("/users/password/reset", useractions.HandlePasswordResetSend).Post()
	r.Add("/users/password/sent", useractions.HandlePasswordResetSentShow)

	// Standard REST handlers for pages
	r.Add("/pages", pageactions.HandleIndex)
	r.Add("/pages/create", pageactions.HandleCreateShow)
	r.Add("/pages/create", pageactions.HandleCreate).Post()
	r.Add("/pages/{id:[0-9]+}/update", pageactions.HandleUpdateShow)
	r.Add("/pages/{id:[0-9]+}/update", pageactions.HandleUpdate).Post()
	r.Add("/pages/{id:[0-9]+}/destroy", pageactions.HandleDestroy).Post()
	r.Add("/pages/{id:[0-9]+}", pageactions.HandleShow)

	// Setup for an empty website
	r.Add("/fragmenta/setup", pageactions.HandleShowSetup)
	r.Add("/fragmenta/setup", pageactions.HandleSetup).Post()

	// Final wildcard route for pages
	r.Add("/{path:[a-z0-9]+}", pageactions.HandleShowPath)

}
