package tagactions

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
	"github.com/fragmenta/fragmenta-cms/src/tags"
)

// names is used to test setting and getting the first string field of the tag.
var names = []string{"foo", "bar"}

// testSetup performs setup for integration tests
// using the test database, real views, and mock authorisation
// If we can run this once for global tests it might be more efficient?
func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(3)
	if err != nil {
		fmt.Printf("tags: Setup db failed %s", err)
	}

	// Set up mock auth
	resource.SetupAuthorisation()

	// Load templates for rendering
	resource.SetupView(3)

	router := mux.New()
	mux.SetDefault(router)

	// FIXME - Need to write routes out here again, but without pkg prefix
	// Any neat way to do this instead? We'd need a separate routes package under app...
	router.Add("/tags", nil)
	router.Add("/tags/create", nil)
	router.Add("/tags/create", nil).Post()
	router.Add("/tags/login", nil)
	router.Add("/tags/login", nil).Post()
	router.Add("/tags/login", nil).Post()
	router.Add("/tags/logout", nil).Post()
	router.Add("/tags/{id:\\d+}/update", nil)
	router.Add("/tags/{id:\\d+}/update", nil).Post()
	router.Add("/tags/{id:\\d+}/destroy", nil).Post()
	router.Add("/tags/{id:\\d+}", nil)

	// Delete all tags to ensure we get consistent results
	query.ExecSQL("delete from tags;")
	query.ExecSQL("ALTER SEQUENCE tags_id_seq RESTART WITH 1;")
}

// Test GET /tags/create
func TestShowCreateTags(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/tags/create", nil)
	w := httptest.NewRecorder()

	// Set up tag session cookie for admin tag above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("tagactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleCreateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Errorf("tagactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Errorf("tagactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /tags/create
func TestCreateTags(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[0])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/tags/create", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up tag session cookie for admin tag
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("tagactions: error setting session %s", err)
	}

	// Run the handler to update the tag
	err = HandleCreate(w, r)
	if err != nil {
		t.Fatalf("tagactions: error handling HandleCreate %s", err)
	}

	// Test we get a redirect after update (to the tag concerned)
	if w.Code != http.StatusFound {
		t.Errorf("tagactions: unexpected response code for HandleCreate expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the tag name is in now value names[1]
	allTags, err := tags.FindAll(tags.Query().Order("id desc"))
	if err != nil || len(allTags) == 0 {
		t.Fatalf("tagactions: error finding created tag %s", err)
	}

	newTags := allTags[0]
	if newTags.ID != 1 || newTags.Name != names[0] {
		t.Errorf("tagactions: error with created tag values: %v %s", newTags.ID, newTags.Name)
	}
}

// Test GET /tags
func TestListTags(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/tags", nil)
	w := httptest.NewRecorder()

	// Set up tag session cookie for admin tag above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("tagactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleIndex(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Errorf("tagactions: error handling HandleIndex %s", err)
	}

	// Test the body for a known pattern
	pattern := "data-table-head"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Errorf("tagactions: unexpected response for HandleIndex expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test of GET /tags/1
func TestShowTags(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/tags/1", nil)
	w := httptest.NewRecorder()

	// Set up tag session cookie for admin tag above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("tagactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Errorf("tagactions: error handling HandleShow %s", err)
	}

	// Test the body for a known pattern
	pattern := names[0]
	if !strings.Contains(w.Body.String(), names[0]) {
		t.Errorf("tagactions: unexpected response for HandleShow expected:%s got:%s", pattern, w.Body.String())
	}
}

// Test GET /tags/123/update
func TestShowUpdateTags(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/tags/1/update", nil)
	w := httptest.NewRecorder()

	// Set up tag session cookie for admin tag above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("tagactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleUpdateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Errorf("tagactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Errorf("tagactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /tags/123/update
func TestUpdateTags(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[1])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/tags/1/update", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up tag session cookie for admin tag
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("tagactions: error setting session %s", err)
	}

	// Run the handler to update the tag
	err = HandleUpdate(w, r)
	if err != nil {
		t.Errorf("tagactions: error handling HandleUpdateTags %s", err)
	}

	// Test we get a redirect after update (to the tag concerned)
	if w.Code != http.StatusFound {
		t.Errorf("tagactions: unexpected response code for HandleUpdateTags expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the tag name is in now value names[1]
	tag, err := tags.Find(1)
	if err != nil {
		t.Fatalf("tagactions: error finding updated tag %s", err)
	}
	if tag.ID != 1 || tag.Name != names[1] {
		t.Errorf("tagactions: error with updated tag values: %v", tag)
	}

}

// Test of POST /tags/123/destroy
func TestDeleteTags(t *testing.T) {

	body := strings.NewReader(``)

	// Now test deleting the tag created above as admin
	r := httptest.NewRequest("POST", "/tags/1/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up tag session cookie for admin tag
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("tagactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleDestroy(w, r)

	// Test the error response is 302 StatusFound
	if err != nil {
		t.Errorf("tagactions: error handling HandleDestroy %s", err)
	}

	// Test we get a redirect after delete
	if w.Code != http.StatusFound {
		t.Errorf("tagactions: unexpected response code for HandleDestroy expected:%d got:%d", http.StatusFound, w.Code)
	}
	// Now test as anon
	r = httptest.NewRequest("POST", "/tags/2/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()

	// Run the handler to test failure as anon
	err = HandleDestroy(w, r)
	if err == nil { // failure expected
		t.Errorf("tagactions: unexpected response for HandleDestroy as anon, expected failure")
	}

}
