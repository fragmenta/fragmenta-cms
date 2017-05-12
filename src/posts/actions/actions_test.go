package postactions

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
	"github.com/fragmenta/fragmenta-cms/src/posts"
)

// names is used to test setting and getting the first string field of the post.
var names = []string{"foo", "bar"}

// testSetup performs setup for integration tests
// using the test database, real views, and mock authorisation
// If we can run this once for global tests it might be more efficient?
func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(3)
	if err != nil {
		fmt.Printf("posts: Setup db failed %s", err)
	}

	// Set up mock auth
	resource.SetupAuthorisation()

	// Load templates for rendering
	resource.SetupView(3)

	router := mux.New()
	mux.SetDefault(router)

	// FIXME - Need to write routes out here again, but without pkg prefix
	// Any neat way to do this instead? We'd need a separate routes package under app...
	router.Add("/posts", nil)
	router.Add("/posts/create", nil)
	router.Add("/posts/create", nil).Post()
	router.Add("/posts/login", nil)
	router.Add("/posts/login", nil).Post()
	router.Add("/posts/login", nil).Post()
	router.Add("/posts/logout", nil).Post()
	router.Add("/posts/{id:\\d+}/update", nil)
	router.Add("/posts/{id:\\d+}/update", nil).Post()
	router.Add("/posts/{id:\\d+}/destroy", nil).Post()
	router.Add("/posts/{id:\\d+}", nil)

	// Delete all posts to ensure we get consistent results
	query.ExecSQL("delete from posts;")
	query.ExecSQL("ALTER SEQUENCE posts_id_seq RESTART WITH 1;")
}

// Test GET /posts/create
func TestShowCreatePost(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/posts/create", nil)
	w := httptest.NewRecorder()

	// Set up post session cookie for admin post above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("postactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleCreateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("postactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("postactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /posts/create
func TestCreatePost(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[0])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/posts/create", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up post session cookie for admin post
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("postactions: error setting session %s", err)
	}

	// Run the handler to update the post
	err = HandleCreate(w, r)
	if err != nil {
		t.Fatalf("postactions: error handling HandleCreate %s", err)
	}

	// Test we get a redirect after update (to the post concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("postactions: unexpected response code for HandleCreate expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the post name is in now value names[1]
	allPost, err := posts.FindAll(posts.Query().Order("id desc"))
	if err != nil || len(allPost) == 0 {
		t.Fatalf("postactions: error finding created post %s", err)
	}
	newPost := allPost[0]
	if newPost.ID != 1 || newPost.Name != names[0] {
		t.Fatalf("postactions: error with created post values: %v %s", newPost.ID, newPost.Name)
	}
}

// Test GET /posts
func TestListPost(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/posts", nil)
	w := httptest.NewRecorder()

	// Set up post session cookie for admin post above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("postactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleIndex(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("postactions: error handling HandleIndex %s", err)
	}

	// Test the body for a known pattern
	pattern := "data-table-head"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("postactions: unexpected response for HandleIndex expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test of GET /posts/1
func TestShowPost(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/posts/1", nil)
	w := httptest.NewRecorder()

	// Set up post session cookie for admin post above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("postactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("postactions: error handling HandleShow %s", err)
	}

	// Test the body for a known pattern
	pattern := names[0]
	if !strings.Contains(w.Body.String(), names[0]) {
		t.Fatalf("postactions: unexpected response for HandleShow expected:%s got:%s", pattern, w.Body.String())
	}
}

// Test GET /posts/123/update
func TestShowUpdatePost(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/posts/1/update", nil)
	w := httptest.NewRecorder()

	// Set up post session cookie for admin post above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("postactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleUpdateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("postactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("postactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /posts/123/update
func TestUpdatePost(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[1])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/posts/1/update", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up post session cookie for admin post
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("postactions: error setting session %s", err)
	}

	// Run the handler to update the post
	err = HandleUpdate(w, r)
	if err != nil {
		t.Fatalf("postactions: error handling HandleUpdatePost %s", err)
	}

	// Test we get a redirect after update (to the post concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("postactions: unexpected response code for HandleUpdatePost expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the post name is in now value names[1]
	post, err := posts.Find(1)
	if err != nil {
		t.Fatalf("postactions: error finding updated post %s", err)
	}
	if post.ID != 1 || post.Name != names[1] {
		t.Fatalf("postactions: error with updated post values: %v", post)
	}

}

// Test of POST /posts/123/destroy
func TestDeletePost(t *testing.T) {

	body := strings.NewReader(``)

	// Now test deleting the post created above as admin
	r := httptest.NewRequest("POST", "/posts/1/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up post session cookie for admin post
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("postactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleDestroy(w, r)

	// Test the error response is 302 StatusFound
	if err != nil {
		t.Fatalf("postactions: error handling HandleDestroy %s", err)
	}

	// Test we get a redirect after delete
	if w.Code != http.StatusFound {
		t.Fatalf("postactions: unexpected response code for HandleDestroy expected:%d got:%d", http.StatusFound, w.Code)
	}
	// Now test as anon
	r = httptest.NewRequest("POST", "/posts/1/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()

	// Run the handler to test failure as anon
	err = HandleDestroy(w, r)
	if err == nil { // failure expected
		t.Fatalf("postactions: unexpected response for HandleDestroy as anon, expected failure")
	}

}
