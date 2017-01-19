package main

import (
	"fmt"
	"net/http"
	"testing"
)

// TestServer tests running the server and using http client to GET /
// this test will fail if there is no secrets file.
func TestServer(t *testing.T) {
	// Setup our server from config
	s, err := SetupServer()
	if err != nil {
		t.Fatalf("server: error setting up %s\n", err)
	}

	// Start the server
	go s.Start()

	// Try hitting the server to see if it is working

	host := fmt.Sprintf("http://localhost%s/", s.PortString())
	r, err := http.Get(host)
	if err != nil {
		t.Fatalf("server: error getting /  %s", err)
	}

	if r.StatusCode != http.StatusOK {
		t.Fatalf("server: error getting / expected:%d got:%d", http.StatusOK, r.StatusCode)
	}

}
