package main

import (
	"fmt"

	"github.com/fragmenta/server"

	"github.com/fragmenta/fragmenta-cms/src/app"
)

func main() {

	// If we have no config, bootstrap first by generating config/migrations
	if app.RequiresBootStrap() {
		err := app.Bootstrap()
		if err != nil {
			fmt.Printf("Error bootstrapping server %s\n", err)
			return
		}
	}

	// Setup server
	server, err := server.New()
	if err != nil {
		fmt.Printf("Error creating server %s\n", err)
		return
	}

	app.Setup(server)

	// Inform user of server setup
	server.Logf("#info TEST Starting server in %s mode on port %d", server.Mode(), server.Port())

	// Start the server
	err = server.Start()
	if err != nil {
		server.Fatalf("Error starting server %s", err)
	}

}
