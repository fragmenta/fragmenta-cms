package app

import (
	"os"
	"time"

	"github.com/fragmenta/assets"
	"github.com/fragmenta/query"
	"github.com/fragmenta/server/config"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/mail"
	"github.com/fragmenta/fragmenta-cms/src/lib/mail/adapters/sendgrid"
)

// appAssets is a pkg global used in our default handlers to serve asset files.
var appAssets *assets.Collection

// Setup sets up our application.
func Setup() {

	// Setup log
	err := SetupLog()
	if err != nil {
		println("failed to set up logs %s", err)
		os.Exit(1)
	}

	// Log server startup
	msg := "Starting server"
	if config.Production() {
		msg = msg + " in production"
	}

	log.Info(log.Values{"msg": msg, "port": config.Get("port")})
	defer log.Time(time.Now(), log.Values{"msg": "Finished loading server"})

	// Set up our mail adapter
	SetupMail()

	// Set up our assets
	SetupAssets()

	// Setup our view templates
	SetupView()

	// Setup our database
	SetupDatabase()

	// Set up auth pkg and authorisation for access
	SetupAuth()

	// Set up our app routes
	SetupRoutes()

}

// SetupLog sets up logging
func SetupLog() error {

	// Set up a stderr logger with time prefix
	logger, err := log.NewStdErr(log.PrefixDateTime)
	if err != nil {
		return err
	}
	log.Add(logger)

	// Set up a file logger pointing at the right location for this config.
	fileLog, err := log.NewFile(config.Get("log"))
	if err != nil {
		return err
	}
	log.Add(fileLog)

	return nil
}

// SetupMail sets us up to send mail via sendgrid (requires key).
func SetupMail() {
	mail.Production = config.Production()
	mail.Service = sendgrid.New(config.Get("mail_from"), config.Get("mail_secret"))
}

// SetupAssets compiles or copies our assets from src into the public assets folder.
func SetupAssets() {
	defer log.Time(time.Now(), log.V{"msg": "Finished loading assets"})

	// Compilation of assets is done on deploy
	// We just load them here
	assetsCompiled := config.GetBool("assets_compiled")

	// Init the pkg global for use in ServeAssets
	appAssets = assets.New(assetsCompiled)

	// Load asset details from json file on each run
	err := appAssets.Load()
	if err != nil {
		// Compile assets for the first time
		log.Info(log.V{"msg": "Compiling Asssets"})

		err := appAssets.Compile("src", "public")
		if err != nil {
			log.Fatal(log.V{"a": "unable to compile assets", "error": err})
			os.Exit(1)
		}
	}

	// Set up helpers which are aware of fingerprinted assets
	// These behave differently depending on the compile flag above
	// when compile is set to no, they use precompiled assets
	// otherwise they serve all files in a group separately
	view.Helpers["style"] = appAssets.StyleLink
	view.Helpers["script"] = appAssets.ScriptLink

}

// SetupView sets up the view package by loadind templates.
func SetupView() {
	defer log.Time(time.Now(), log.V{"msg": "Finished loading templates"})

	view.Production = config.Production()
	err := view.LoadTemplates()
	if err != nil {
		//	server.Fatalf("Error reading templates %s", err)
		log.Fatal(log.V{"msg": "unable to read templates", "error": err})
		os.Exit(1)
	}

}

// SetupDatabase sets up the db with query given our server config.
func SetupDatabase() {
	defer log.Time(time.Now(), log.V{"msg": "Finished opening database", "db": config.Get("db"), "user": config.Get("db_user")})

	options := map[string]string{
		"adapter":  config.Get("db_adapter"),
		"user":     config.Get("db_user"),
		"password": config.Get("db_pass"),
		"db":       config.Get("db"),
	}

	// Ask query to open the database
	err := query.OpenDatabase(options)

	if err != nil {
		log.Fatal(log.V{"msg": "unable to read database", "db": config.Get("db"), "error": err})
		os.Exit(1)
	}

}
