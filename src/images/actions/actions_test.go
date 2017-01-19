package imageactions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fragmenta/mux"
	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-cms/src/images"
	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
)

// names is used to test setting and getting the first string field of the image.
var names = []string{"foo", "bar"}

// testSetup performs setup for integration tests
// using the test database, real views, and mock authorisation
// If we can run this once for global tests it might be more efficient?
func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(3)
	if err != nil {
		fmt.Printf("images: Setup db failed %s", err)
	}

	// Set up mock auth
	resource.SetupAuthorisation()

	// Load templates for rendering
	resource.SetupView(3)

	router := mux.New()
	mux.SetDefault(router)

	// FIXME - Need to write routes out here again, but without pkg prefix
	// Any neat way to do this instead? We'd need a separate routes package under app...
	router.Add("/images", nil)
	router.Add("/images/create", nil)
	router.Add("/images/create", nil).Post()
	router.Add("/images/login", nil)
	router.Add("/images/login", nil).Post()
	router.Add("/images/login", nil).Post()
	router.Add("/images/logout", nil).Post()
	router.Add("/images/{id:\\d+}/update", nil)
	router.Add("/images/{id:\\d+}/update", nil).Post()
	router.Add("/images/{id:\\d+}/destroy", nil).Post()
	router.Add("/images/{id:\\d+}", nil)

	// Delete all images to ensure we get consistent results
	query.ExecSQL("delete from images;")
	query.ExecSQL("ALTER SEQUENCE images_id_seq RESTART WITH 1;")
}

// Test GET /images/create
func TestShowCreateImages(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/images/create", nil)
	w := httptest.NewRecorder()

	// Set up image session cookie for admin image above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("imageactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleCreateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("imageactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("imageactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /images/create
func TestCreateImages(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[0])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/images/create", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up image session cookie for admin image
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("imageactions: error setting session %s", err)
	}

	// Run the handler to update the image
	err = HandleCreate(w, r)
	if err != nil {
		t.Fatalf("imageactions: error handling HandleCreate %s", err)
	}

	// Test we get a redirect after update (to the image concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("imageactions: unexpected response code for HandleCreate expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the image name is in now value names[1]
	allImages, err := images.FindAll(images.Query().Order("id desc"))
	if err != nil || len(allImages) == 0 {
		t.Fatalf("imageactions: error finding created image %s", err)
	}
	newImages := allImages[0]
	if newImages.ID != 1 || newImages.Name != names[0] {
		t.Fatalf("imageactions: error with created image values: %v %s", newImages.ID, newImages.Name)
	}
}

// Test GET /images
func TestListImages(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/images", nil)
	w := httptest.NewRecorder()

	// Set up image session cookie for admin image above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("imageactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleIndex(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("imageactions: error handling HandleIndex %s", err)
	}

	// Test the body for a known pattern
	pattern := "data-table-head"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("imageactions: unexpected response for HandleIndex expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test of GET /images/1
func TestShowImages(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/images/1", nil)
	w := httptest.NewRecorder()

	// Set up image session cookie for admin image above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("imageactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("imageactions: error handling HandleShow %s", err)
	}

	// Test the body for a known pattern
	pattern := names[0]
	if !strings.Contains(w.Body.String(), names[0]) {
		t.Fatalf("imageactions: unexpected response for HandleShow expected:%s got:%s", pattern, w.Body.String())
	}
}

// Test GET /images/123/update
func TestShowUpdateImages(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/images/1/update", nil)
	w := httptest.NewRecorder()

	// Set up image session cookie for admin image above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("imageactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleUpdateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("imageactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("imageactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /images/123/update
func TestUpdateImages(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[1])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/images/1/update", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up image session cookie for admin image
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("imageactions: error setting session %s", err)
	}

	// Run the handler to update the image
	err = HandleUpdate(w, r)
	if err != nil {
		t.Fatalf("imageactions: error handling HandleUpdateImages %s", err)
	}

	// Test we get a redirect after update (to the image concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("imageactions: unexpected response code for HandleUpdateImages expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the image name is in now value names[1]
	image, err := images.Find(1)
	if err != nil {
		t.Fatalf("imageactions: error finding updated image %s", err)
	}
	if image.ID != 1 || image.Name != names[1] {
		t.Fatalf("imageactions: error with updated image values: %v", image)
	}

}

// Test of POST /images/123/destroy
func TestDeleteImages(t *testing.T) {

	body := strings.NewReader(``)

	// Now test deleting the image created above as admin
	r := httptest.NewRequest("POST", "/images/1/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up image session cookie for admin image
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("imageactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleDestroy(w, r)

	// Test the error response is 302 StatusFound
	if err != nil {
		t.Fatalf("imageactions: error handling HandleDestroy %s", err)
	}

	// Test we get a redirect after delete
	if w.Code != http.StatusFound {
		t.Fatalf("imageactions: unexpected response code for HandleDestroy expected:%d got:%d", http.StatusFound, w.Code)
	}
	// Now test as anon
	r = httptest.NewRequest("POST", "/images/1/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()

	// Run the handler to test failure as anon
	err = HandleDestroy(w, r)
	if err == nil { // failure expected
		t.Fatalf("imageactions: unexpected response for HandleDestroy as anon, expected failure")
	}

}
