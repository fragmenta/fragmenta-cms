package main

import (
	"fmt"

	"github.com/fragmenta/server"

	"github.com/fragmenta/fragmenta-cms/src/app"
)

func main() {

	// Setup server
	server, err := server.New()
	if err != nil {
		fmt.Printf("Error creating server %s", err)
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
