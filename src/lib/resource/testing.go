package resource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/query"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"
)

// This file contains some test helpers for resources.

// basePath returns the path to the fragmenta root from a given test folder.
func basePath(depth int) string {
	// Construct a path to root
	p := ""
	for i := 0; i < depth; i++ {
		p = filepath.Join(p, "..")
	}
	return p
}

// SetupAuthorisation sets up mock authorisation.
func SetupAuthorisation() {

	// Setup the auth library with fake data
	var testKey = "12353bce2bbc4efb90eff81c29dc982de9a0176b568db18a61b4f4732cadabbc"

	// Setup auth with some test values - could read these from config I guess
	auth.HMACKey = auth.HexToBytes(testKey)
	auth.SecretKey = auth.HexToBytes(testKey)
	auth.SessionName = "test_session"

	// Set up admin permissions for testing -
	// hard coded role to avoid cyclic dependency
	can.Authorise(100, can.ManageResource, can.Anything)
}

// AddUserSessionCookie adds a new cookie for the given user
// on the incoming request, so that we can test authentication in handlers.
func AddUserSessionCookie(w *httptest.ResponseRecorder, r *http.Request, id int) error {

	// Build the session from the secure cookie, or create a new one
	session, err := auth.Session(w, r)
	if err != nil {
		return err
	}

	secret := auth.BytesToBase64(auth.RandomToken(auth.TokenLength))
	session.Set(auth.SessionTokenKey, secret)

	// Now from secret, generate a secure token for this request
	token := auth.BytesToBase64(auth.AuthenticityTokenWithSecret(auth.Base64ToBytes(secret)))

	// Write value of user id 1
	session.Set(auth.SessionUserKey, "1")

	// Set the cookie on the recorder
	err = session.Save(w)
	if err != nil {
		return err
	}

	// Set the auth token on params of request
	// Cheat and set on raw query, which we don't use in tests
	urlQ := fmt.Sprintf("authenticity_token=%s", token)
	r.URL.RawQuery = urlQ

	// Now get the entire cookie back out
	// and put it on the request as if it were coming in from browser
	r.Header.Set("Cookie", strings.Join(w.HeaderMap["Set-Cookie"], ""))

	// Perform an authenticity check:
	err = auth.CheckAuthenticityToken(token, r)
	if err != nil {
		return err
	}

	return nil
}

// SetupView sets up the view package for testing by loading templates.
func SetupView(depth int) error {
	view.Production = false
	return view.LoadTemplatesAtPaths([]string{filepath.Join(basePath(depth), "src")}, view.Helpers)
}

// SetupTestDatabase sets up the database for all tests from the test config.
func SetupTestDatabase(depth int) error {

	// Set up a stderr logger with time prefix
	logger, err := log.NewStdErr(log.PrefixDateTime)
	if err != nil {
		return err
	}
	log.Add(logger)

	// Read config json
	path := filepath.Join(basePath(depth), "secrets", "fragmenta.json")
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var data map[string]map[string]string
	err = json.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	config := data["test"]
	options := map[string]string{
		"adapter":  config["db_adapter"],
		"user":     config["db_user"],
		"password": config["db_pass"],
		"db":       config["db"],
	}

	// Ask query to open the database
	err = query.OpenDatabase(options)
	if err != nil {
		return err
	}

	// For speed
	query.Exec("set synchronous_commit=off;")
	return nil
}
