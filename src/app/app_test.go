package app

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/server/config"

	"github.com/fragmenta/fragmenta-cms/src/users"
)

// TestRouter tests our routes are functioning correctly.
func TestRouter(t *testing.T) {

	// chdir into the root to load config/assets to test this code
	err := os.Chdir("../../")
	if err != nil {
		t.Errorf("Chdir error: %s", err)
	}

	c := config.New()
	c.Load("secrets/fragmenta.json")
	c.Mode = config.ModeTest
	config.Current = c

	// First, set up the logger
	err = SetupLog()
	if err != nil {
		t.Fatalf("app: failed to set up log %s", err)
	}

	// Set up our assets
	SetupAssets()

	// Setup our view templates
	SetupView()

	// Setup our database
	SetupDatabase()

	// Set up auth pkg and authorisation for access
	SetupAuth()

	// Setup our router and handlers
	router := SetupRoutes()

	// Test serving the route / which should always exist
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	// Test code on response
	if w.Code != http.StatusOK {
		t.Fatalf("app: error code on / expected:%d got:%d", http.StatusOK, w.Code)
	}

}

// TestAuth tests our authentication is functioning after setup.
func TestAuth(t *testing.T) {

	SetupAuth()

	user := &users.User{}

	// Test anon cannot access /users
	err := can.List(user, users.MockAnon())
	if err == nil {
		t.Fatalf("app: authentication block failed for anon")
	}

	// Test anon cannot edit admin user
	err = can.Update(users.MockAdmin(), users.MockAnon())
	if err == nil {
		t.Fatalf("app: authentication block failed for anon")
	}

	// Test admin can access /users
	err = can.List(user, users.MockAdmin())
	if err != nil {
		t.Fatalf("app: authentication failed for admin")
	}

	// Test admin can admin user
	err = can.Manage(user, users.MockAdmin())
	if err != nil {
		t.Fatalf("app: authentication failed for admin")
	}

}
