package main

import (
	"fmt"
	"os"

	"github.com/fragmenta/server"
	"github.com/fragmenta/server/config"

	"github.com/fragmenta/fragmenta-cms/src/app"
)

// Main entrypoint for the server which performs bootstrap, setup
// then runs the server. Most setup is delegated to the src/app pkg.
func main() {

	// Bootstrap if required (no config file found).
	if app.RequiresBootStrap() {
		err := app.Bootstrap()
		if err != nil {
			fmt.Printf("Error bootstrapping server %s\n", err)
			return
		}
	}

	// Setup our server
	s, err := SetupServer()
	if err != nil {
		fmt.Printf("server: error setting up %s\n", err)
		return
	}

	// Start the server
	err = s.Start()
	if err != nil {
		s.Fatalf("server: error starting %s\n", err)
	}

}

// SetupServer creates a new server, and delegates setup to the app pkg.
func SetupServer() (*server.Server, error) {

	// Setup server
	s, err := server.New()
	if err != nil {
		return nil, err
	}

	// Load the appropriate config
	c := config.New()
	err = c.Load("secrets/fragmenta.json")
	if err != nil {
		return nil, err
	}
	config.Current = c

	// Check environment variable to see if we are in production mode
	if os.Getenv("FRAG_ENV") == "production" {
		config.Current.Mode = config.ModeProduction
	}

	// Call the app to perform additional setup
	app.Setup()

	return s, nil
}
