package session

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fragmenta/auth"
)

var (
	testKey = "12353bce2bbc4efb90eff81c29dc982de9a0176b568db18a61b4f4732cadabbc"
	set     = "foo"
)

// TestAuthenticate tests storing a value in a cookie and retreiving it again.
func TestAuthenticate(t *testing.T) {
	// Setup auth with some test values - could read these from config I guess
	auth.HMACKey = auth.HexToBytes(testKey)
	auth.SecretKey = auth.HexToBytes(testKey)
	auth.SessionName = "test_session"

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Get the current user for this request (nil)
	currentUser := CurrentUser(w, r)
	if currentUser.ID != 0 {
		t.Fatalf("auth: failed to get empty user")
	}

	// Now set a user on session on request, and try again
	session, err := auth.Session(w, r)
	if err != nil {
		t.Fatalf("auth: failed to build session")
	}
	session.Set(auth.SessionUserKey, "1")

	// Set the cookie on the recorder
	err = session.Save(w)
	if err != nil {
		t.Fatalf("auth: failed to save session")
	}

	// Now get the cookie back out and put it on the request as if it were coming in from browser
	r.Header.Set("Cookie", strings.Join(w.HeaderMap["Set-Cookie"], ""))

	t.Logf("SESSION:%v", session)

	session, err = auth.Session(w, r)
	if err != nil {
		t.Fatalf("auth: failed to build session")
	}
	if session.Get(auth.SessionUserKey) != "1" {
		t.Fatalf("auth: failed to restore session value")
	}

}
