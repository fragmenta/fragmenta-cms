package pageactions

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
	"github.com/fragmenta/fragmenta-cms/src/pages"
)

// names is used to test setting and getting the first string field of the page.
var names = []string{"foo", "bar"}

// testSetup performs setup for integration tests
// using the test database, real views, and mock authorisation
// If we can run this once for global tests it might be more efficient?
func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(3)
	if err != nil {
		fmt.Printf("pages: Setup db failed %s", err)
	}

	// Set up mock auth
	resource.SetupAuthorisation()

	// Load templates for rendering
	resource.SetupView(3)

	router := mux.New()
	mux.SetDefault(router)

	// FIXME - Need to write routes out here again, but without pkg prefix
	// Any neat way to do this instead? We'd need a separate routes package under app...
	router.Add("/pages", nil)
	router.Add("/pages/create", nil)
	router.Add("/pages/create", nil).Post()
	router.Add("/pages/login", nil)
	router.Add("/pages/login", nil).Post()
	router.Add("/pages/login", nil).Post()
	router.Add("/pages/logout", nil).Post()
	router.Add("/pages/{id:\\d+}/update", nil)
	router.Add("/pages/{id:\\d+}/update", nil).Post()
	router.Add("/pages/{id:\\d+}/destroy", nil).Post()
	router.Add("/pages/{id:\\d+}", nil)

	// Delete all pages to ensure we get consistent results
	query.ExecSQL("delete from pages;")
	query.ExecSQL("ALTER SEQUENCE pages_id_seq RESTART WITH 1;")
}

// Test GET /pages/create
func TestShowCreatePage(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/pages/create", nil)
	w := httptest.NewRecorder()

	// Set up page session cookie for admin page above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("pageactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleCreateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("pageactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("pageactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /pages/create
func TestCreatePage(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[0])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/pages/create", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up page session cookie for admin page
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("pageactions: error setting session %s", err)
	}

	// Run the handler to update the page
	err = HandleCreate(w, r)
	if err != nil {
		t.Fatalf("pageactions: error handling HandleCreate %s", err)
	}

	// Test we get a redirect after update (to the page concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("pageactions: unexpected response code for HandleCreate expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the page name is in now value names[1]
	allPage, err := pages.FindAll(pages.Query().Order("id desc"))
	if err != nil || len(allPage) == 0 {
		t.Fatalf("pageactions: error finding created page %s", err)
	}
	newPage := allPage[0]
	if newPage.ID != 1 || newPage.Name != names[0] {
		t.Fatalf("pageactions: error with created page values: %v %s", newPage.ID, newPage.Name)
	}
}

// Test GET /pages
func TestListPage(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/pages", nil)
	w := httptest.NewRecorder()

	// Set up page session cookie for admin page above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("pageactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleIndex(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("pageactions: error handling HandleIndex %s", err)
	}

	// Test the body for a known pattern
	pattern := "data-table-head"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("pageactions: unexpected response for HandleIndex expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test of GET /pages/1
func TestShowPage(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/pages/1", nil)
	w := httptest.NewRecorder()

	// Set up page session cookie for admin page above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("pageactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("pageactions: error handling HandleShow %s", err)
	}

	// Test the body for a known pattern
	pattern := names[0]
	if !strings.Contains(w.Body.String(), names[0]) {
		t.Fatalf("pageactions: unexpected response for HandleShow expected:%s got:%s", pattern, w.Body.String())
	}
}

// Test GET /pages/123/update
func TestShowUpdatePage(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/pages/1/update", nil)
	w := httptest.NewRecorder()

	// Set up page session cookie for admin page above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("pageactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleUpdateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("pageactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("pageactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /pages/123/update
func TestUpdatePage(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[1])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/pages/1/update", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up page session cookie for admin page
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("pageactions: error setting session %s", err)
	}

	// Run the handler to update the page
	err = HandleUpdate(w, r)
	if err != nil {
		t.Fatalf("pageactions: error handling HandleUpdatePage %s", err)
	}

	// Test we get a redirect after update (to the page concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("pageactions: unexpected response code for HandleUpdatePage expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the page name is in now value names[1]
	page, err := pages.Find(1)
	if err != nil {
		t.Fatalf("pageactions: error finding updated page %s", err)
	}
	if page.ID != 1 || page.Name != names[1] {
		t.Fatalf("pageactions: error with updated page values: %v", page)
	}

}

// Test of POST /pages/123/destroy
func TestDeletePage(t *testing.T) {

	body := strings.NewReader(``)

	// Now test deleting the page created above as admin
	r := httptest.NewRequest("POST", "/pages/1/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up page session cookie for admin page
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("pageactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleDestroy(w, r)

	// Test the error response is 302 StatusFound
	if err != nil {
		t.Fatalf("pageactions: error handling HandleDestroy %s", err)
	}

	// Test we get a redirect after delete
	if w.Code != http.StatusFound {
		t.Fatalf("pageactions: unexpected response code for HandleDestroy expected:%d got:%d", http.StatusFound, w.Code)
	}
	// Now test as anon
	r = httptest.NewRequest("POST", "/pages/1/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()

	// Run the handler to test failure as anon
	err = HandleDestroy(w, r)
	if err == nil { // failure expected
		t.Fatalf("pageactions: unexpected response for HandleDestroy as anon, expected failure")
	}

}
