package useractions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

// names is used to test setting and getting the first string field of the user.
var names = []string{"foo", "bar"}

// testSetup performs setup for integration tests
// using the test database, real views, and mock authorisation
// If we can run this once for global tests it might be more efficient?
func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(3)
	if err != nil {
		fmt.Printf("users: Setup db failed %s", err)
	}

	// Set up mock auth
	resource.SetupAuthorisation()

	// Load templates for rendering
	resource.SetupView(3)

	router := mux.New()
	mux.SetDefault(router)

	// FIXME - would prefer to load real routes from app pkg
	router.Add("/users", nil)
	router.Add("/users/create", nil)
	router.Add("/users/create", nil).Post()
	router.Add("/users/login", nil)
	router.Add("/users/login", nil).Post()
	router.Add("/users/login", nil).Post()
	router.Add("/users/logout", nil).Post()
	router.Add("/users/{id:\\d+}/update", nil)
	router.Add("/users/{id:\\d+}/update", nil).Post()
	router.Add("/users/{id:\\d+}/destroy", nil).Post()
	router.Add("/users/{id:\\d+}", nil)

	// Delete all users to ensure we get consistent results?
	_, err = query.ExecSQL("delete from users;")
	if err != nil {
		t.Fatalf("error setting up:%s", err)
	}
	// Insert a test admin user for checking logins - never delete as will
	// be required for other resources testing
	_, err = query.ExecSQL("INSERT INTO users (id,email,name,status,role,password_hash) VALUES(1,'example@example.com','test',100,100,'$2a$10$2IUzpI/yH0Xc.qs9Z5UUL.3f9bqi0ThvbKs6Q91UOlyCEGY8hdBw6');")
	if err != nil {
		t.Fatalf("error setting up:%s", err)
	}
	// Insert user to delete
	_, err = query.ExecSQL("INSERT INTO users (id,email,name,status,role,password_hash) VALUES(2,'example@example.com','test',100,0,'$2a$10$2IUzpI/yH0Xc.qs9Z5UUL.3f9bqi0ThvbKs6Q91UOlyCEGY8hdBw6');")
	if err != nil {
		t.Fatalf("error setting up:%s", err)
	}
	_, err = query.ExecSQL("ALTER SEQUENCE users_id_seq RESTART WITH 3;")
	if err != nil {
		t.Fatalf("error setting up:%s", err)
	}
}

// Test GET /users/create
func TestShowCreateUser(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/users/create", nil)
	w := httptest.NewRecorder()

	// Set up user session cookie for admin user above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("useractions: error setting session %s", err)
	}

	// Run the handler
	err = HandleCreateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Errorf("useractions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Errorf("useractions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /users/create
func TestCreateUser(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[0])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/users/create", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up user session cookie for admin user
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("useractions: error setting session %s", err)
	}

	// Run the handler to update the user
	err = HandleCreate(w, r)
	if err != nil {
		t.Errorf("useractions: error handling HandleCreate %s", err)
	}

	// Test we get a redirect after update (to the user concerned)
	if w.Code != http.StatusFound {
		t.Errorf("useractions: unexpected response code for HandleCreate expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the user name is in now value names[1]
	allUsers, err := users.FindAll(users.Query().Order("id desc"))
	if err != nil || len(allUsers) == 0 {
		t.Fatalf("useractions: error finding created user %s", err)
	}
	newUser := allUsers[0]
	if newUser.ID < 2 || newUser.Name != names[0] {
		t.Errorf("useractions: error with created user values: %v %s", newUser.ID, newUser.Name)
	}
}

// Test GET /users
func TestListUsers(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	// Set up user session cookie for admin user above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("useractions: error setting session %s", err)
	}

	// Run the handler
	err = HandleIndex(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Errorf("useractions: error handling HandleIndex %s", err)
	}

	// Test the body for a known pattern
	pattern := "data-table-head"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Errorf("useractions: unexpected response for HandleIndex expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test of GET /users/1
func TestShowUser(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()

	// Set up user session cookie for admin user above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("useractions: error setting session %s", err)
	}

	// Run the handler
	err = HandleShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Errorf("useractions: error handling HandleShow %s", err)
	}

	// Test the body for a known pattern
	pattern := names[0]
	if !strings.Contains(w.Body.String(), names[0]) {
		t.Errorf("useractions: unexpected response for HandleShow expected:%s got:%s", pattern, w.Body.String())
	}
}

// Test GET /users/123/update
func TestShowUpdateUser(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/users/1/update", nil)
	w := httptest.NewRecorder()

	// Set up user session cookie for admin user above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("useractions: error setting session %s", err)
	}

	// Run the handler
	err = HandleUpdateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Errorf("useractions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Errorf("useractions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /users/123/update
func TestUpdateUser(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[1])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/users/1/update", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up user session cookie for admin user
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("useractions: error setting session %s", err)
	}

	// Run the handler to update the user
	err = HandleUpdate(w, r)
	if err != nil {
		t.Errorf("useractions: error handling HandleUpdateUser %s", err)
	}

	// Test we get a redirect after update (to the user concerned)
	if w.Code != http.StatusFound {
		t.Errorf("useractions: unexpected response code for HandleUpdateUser expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the user name is in now value names[1]
	user, err := users.Find(1)
	if err != nil {
		t.Fatalf("useractions: error finding updated user %s", err)
	}
	if user.ID != 1 || user.Name != names[1] {
		t.Errorf("useractions: error with updated user values: %v", user)
	}

}

// Test of POST /users/123/destroy
func TestDeleteUser(t *testing.T) {

	body := strings.NewReader(``)

	// Now test deleting the user created above as admin
	r := httptest.NewRequest("POST", "/users/2/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up user session cookie for admin user
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("useractions: error setting session %s", err)
	}

	// Run the handler
	err = HandleDestroy(w, r)

	// Test the error response is 302 StatusFound
	if err != nil {
		t.Errorf("useractions: error handling HandleDestroy %s", err)
	}

	// Test we get a redirect after delete
	if w.Code != http.StatusFound {
		t.Errorf("useractions: unexpected response code for HandleDestroy expected:%d got:%d", http.StatusFound, w.Code)
	}
	// Now test as anon
	r = httptest.NewRequest("POST", "/users/2/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()

	// Run the handler to test failure as anon
	err = HandleDestroy(w, r)
	if err == nil { // failure expected
		t.Errorf("useractions: unexpected response for HandleDestroy as anon, expected failure")
	}

}

// Test GET /users/login
func TestShowLogin(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/users/login", nil)
	w := httptest.NewRecorder()

	// Set up user session cookie for admin user above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("useractions: error setting session %s", err)
	}

	// Run the handler
	err = HandleLoginShow(w, r)

	// Check for redirect as they are considered logged in
	if err != nil || w.Code != http.StatusFound {
		t.Errorf("useractions: error handling HandleLoginShow %s %d", err, w.Code)
	}

	// Setup new request and recorder with no session
	r = httptest.NewRequest("GET", "/users/login", nil)
	w = httptest.NewRecorder()

	// Run the handler
	err = HandleLoginShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Errorf("useractions: error handling HandleLoginShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "password"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Errorf("useractions: unexpected response for HandleLoginShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /users/login
func TestLogin(t *testing.T) {

	// These need to match entries in the test db for this to work
	form := url.Values{}
	form.Add("email", "example@example.com")
	form.Add("password", "Hunter2")
	body := strings.NewReader(form.Encode())

	// Test posting to the login link,
	// we expect success as setup inserts this user
	r := httptest.NewRequest("POST", "/users/login", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up user session cookie for anon user (for the CSRF cookie token)
	err := resource.AddUserSessionCookie(w, r, 0)
	if err != nil {
		t.Errorf("useractions: error setting session %s", err)
	}

	// Run the handler
	err = HandleLogin(w, r)
	if err != nil || w.Code != http.StatusFound {
		t.Errorf("useractions: error on HandleLogin %s", err)
	}

}

// Test POST /users/logout
func TestLogout(t *testing.T) {

	r := httptest.NewRequest("POST", "/users/logout", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up user session cookie for admin user
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Errorf("useractions: error setting session %s", err)
	}

	// Run the handler
	err = HandleLogout(w, r)
	if err != nil {
		t.Errorf("useractions: error on HandleLogout %s", err)
	}

	t.Logf("SESSION CLEAR: %s", w.Header().Get("Set-Cookie"))

	t.Logf("SESSION AFTER: %s", w.Header())

	// Check we've set an empty session on this outgoing writer
	if !strings.Contains(string(w.Header().Get("Set-Cookie")), auth.SessionName+"=;") {
		t.Errorf("useractions: error on HandleLogout - session not cleared")
	}

	// TODO - to better test this we should have an integration test with a server

}
