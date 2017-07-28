package redirectactions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fragmenta/mux"
	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
	"github.com/fragmenta/fragmenta-cms/src/redirects"
)

// names is used to test setting and getting the first string field of the redirect.
var names = []string{"foo", "bar"}

// testSetup performs setup for integration tests
// using the test database, real views, and mock authorisation
// If we can run this once for global tests it might be more efficient?
func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(3)
	if err != nil {
		fmt.Printf("redirects: Setup db failed %s", err)
	}

	// Set up mock auth
	resource.SetupAuthorisation()

	// Load templates for rendering
	resource.SetupView(3)

	router := mux.New()
	mux.SetDefault(router)

	// FIXME - Need to write routes out here again, but without pkg prefix
	// Any neat way to do this instead? We'd need a separate routes package under app...
	router.Add("/redirects", nil)
	router.Add("/redirects/create", nil)
	router.Add("/redirects/create", nil).Post()
	router.Add("/redirects/login", nil)
	router.Add("/redirects/login", nil).Post()
	router.Add("/redirects/login", nil).Post()
	router.Add("/redirects/logout", nil).Post()
	router.Add("/redirects/{id:\\d+}/update", nil)
	router.Add("/redirects/{id:\\d+}/update", nil).Post()
	router.Add("/redirects/{id:\\d+}/destroy", nil).Post()
	router.Add("/redirects/{id:\\d+}", nil)

	// Delete all redirects to ensure we get consistent results
	query.ExecSQL("delete from redirects;")
	query.ExecSQL("ALTER SEQUENCE redirects_id_seq RESTART WITH 1;")
}

// Test GET /redirects/create
func TestShowCreateRedirect(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/redirects/create", nil)
	w := httptest.NewRecorder()

	// Set up redirect session cookie for admin redirect above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("redirectactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleCreateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("redirectactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("redirectactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /redirects/create
func TestCreateRedirect(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[0])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/redirects/create", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up redirect session cookie for admin redirect
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("redirectactions: error setting session %s", err)
	}

	// Run the handler to update the redirect
	err = HandleCreate(w, r)
	if err != nil {
		t.Fatalf("redirectactions: error handling HandleCreate %s", err)
	}

	// Test we get a redirect after update (to the redirect concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("redirectactions: unexpected response code for HandleCreate expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the redirect name is in now value names[1]
	allRedirects, err := redirects.FindAll(redirects.Query().Order("id desc"))
	if err != nil || len(allRedirects) == 0 {
		t.Fatalf("redirectactions: error finding created redirect %s", err)
	}
	newRedirect := allRedirects[0]
	if newRedirect.ID != 1 || newRedirect.Name != names[0] {
		t.Fatalf("redirectactions: error with created redirect values: %v %s", newRedirect.ID, newRedirect.Name)
	}
}

// Test GET /redirects
func TestListRedirects(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/redirects", nil)
	w := httptest.NewRecorder()

	// Set up redirect session cookie for admin redirect above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("redirectactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleIndex(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("redirectactions: error handling HandleIndex %s", err)
	}

	// Test the body for a known pattern
	pattern := "data-table-head"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("redirectactions: unexpected response for HandleIndex expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test of GET /redirects/1
func TestShowRedirect(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/redirects/1", nil)
	w := httptest.NewRecorder()

	// Set up redirect session cookie for admin redirect above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("redirectactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("redirectactions: error handling HandleShow %s", err)
	}

	// Test the body for a known pattern
	pattern := names[0]
	if !strings.Contains(w.Body.String(), names[0]) {
		t.Fatalf("redirectactions: unexpected response for HandleShow expected:%s got:%s", pattern, w.Body.String())
	}
}

// Test GET /redirects/123/update
func TestShowUpdateRedirect(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/redirects/1/update", nil)
	w := httptest.NewRecorder()

	// Set up redirect session cookie for admin redirect above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("redirectactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleUpdateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("redirectactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("redirectactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /redirects/123/update
func TestUpdateRedirect(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[1])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/redirects/1/update", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up redirect session cookie for admin redirect
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("redirectactions: error setting session %s", err)
	}

	// Run the handler to update the redirect
	err = HandleUpdate(w, r)
	if err != nil {
		t.Fatalf("redirectactions: error handling HandleUpdateRedirect %s", err)
	}

	// Test we get a redirect after update (to the redirect concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("redirectactions: unexpected response code for HandleUpdateRedirect expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the redirect name is in now value names[1]
	redirect, err := redirects.Find(1)
	if err != nil {
		t.Fatalf("redirectactions: error finding updated redirect %s", err)
	}
	if redirect.ID != 1 || redirect.Name != names[1] {
		t.Fatalf("redirectactions: error with updated redirect values: %v", redirect)
	}

}

// Test of POST /redirects/123/destroy
func TestDeleteRedirect(t *testing.T) {

	body := strings.NewReader(``)

	// Now test deleting the redirect created above as admin
	r := httptest.NewRequest("POST", "/redirects/1/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up redirect session cookie for admin redirect
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("redirectactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleDestroy(w, r)

	// Test the error response is 302 StatusFound
	if err != nil {
		t.Fatalf("redirectactions: error handling HandleDestroy %s", err)
	}

	// Test we get a redirect after delete
	if w.Code != http.StatusFound {
		t.Fatalf("redirectactions: unexpected response code for HandleDestroy expected:%d got:%d", http.StatusFound, w.Code)
	}
	// Now test as anon
	r = httptest.NewRequest("POST", "/redirects/1/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()

	// Run the handler to test failure as anon
	err = HandleDestroy(w, r)
	if err == nil { // failure expected
		t.Fatalf("redirectactions: unexpected response for HandleDestroy as anon, expected failure")
	}

}
