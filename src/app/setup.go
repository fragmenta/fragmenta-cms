package app

import (
	"time"

	"github.com/fragmenta/assets"
	"github.com/fragmenta/query"
	"github.com/fragmenta/router"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-app/src/lib/authorise"
)

var appAssets *assets.Collection

func Setup(server *server.Server) {

	// Setup log
	server.Logger = log.New(server.Config("log"), server.Production())

	// Set up our assets
	setupAssets(server)

	// Setup our view templates
	setupView(server)

	// Setup our database
	setupDatabase(server)

	// Routing
	router, err := router.New(server.Logger, server)
	if err != nil {
		server.Fatalf("Error creating router %s", err)
	}

	// Setup our authentication and authorisation
	authorise.Setup(server)

	// Setup our router and handlers
	setupRoutes(router)

}

// Compile or copy in our assets from src into the public assets folder, for use by the app
func setupAssets(server *server.Server) {
	defer server.Timef("#info Finished loading assets in %s", time.Now())

	// Compilation of assets is done on deploy
	// We just load them here
	assetsCompiled := server.ConfigBool("assets_compiled")
	appAssets = assets.New(assetsCompiled)

	// Load asset details from json file on each run
	err := appAssets.Load()
	if err != nil {
		// Compile assets for the first time
		server.Logf("#info Compiling assets")
		err := appAssets.Compile("src", "public")
		if err != nil {
			server.Fatalf("#error compiling assets %s", err)
		}
	}

	// Set up helpers which are aware of fingerprinted assets
	// These behave differently depending on the compile flag above
	// when compile is set to no, they use precompiled assets
	// otherwise they serve all files in a group separately
	view.Helpers["style"] = appAssets.StyleLink
	view.Helpers["script"] = appAssets.ScriptLink

}

func setupView(server *server.Server) {
	// Set up our own func map here if required
	//funcs := view.DefaultHelpers
	//view.Helpers = funcs
	view.Production = server.Production()

	started := time.Now()
	err := view.LoadTemplates()

	if err != nil {
		server.Fatalf("Error reading templates %s", err)
	}

	// Log time taken loading templates
	end := time.Since(started).String()
	server.Logf("#info Parsed templates in %s", end)

}

// Setup db - at present query pkg manages this...
func setupDatabase(server *server.Server) {
	config := server.Configuration()
	options := map[string]string{
		"adapter":  config["db_adapter"],
		"user":     config["db_user"],
		"password": config["db_pass"],
		"db":       config["db"],
	}

	// Ask query to open the database
	err := query.OpenDatabase(options)

	if err != nil {
		server.Fatalf("Error reading database %s", err)
	}

	server.Logf("#info Opened database at %s for user %s", server.Config("db"), server.Config("db_user"))

}
