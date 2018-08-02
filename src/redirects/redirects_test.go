// Tests for the redirects package
package redirects

import (
	"testing"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
)

var testURL = "foo"

func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(2)
	if err != nil {
		t.Fatalf("redirects: Setup db failed %s", err)
	}
}

// Test Create method
func TestCreateRedirect(t *testing.T) {
	redirectParams := map[string]string{
		"new_url": testURL,
	}

	id, err := New().Create(redirectParams)
	if err != nil {
		t.Fatalf("redirects: Create redirect failed :%s", err)
	}

	redirect, err := Find(id)
	if err != nil {
		t.Fatalf("redirects: Create redirect find failed")
	}

	if redirect.NewURL != testURL {
		t.Fatalf("redirects: Create redirect url failed expected:%s got:%s", testURL, redirect.NewURL)
	}

}

// Test Index (List) method
func TestListRedirects(t *testing.T) {

	// Get all redirects (we should have at least one)
	results, err := FindAll(Query())
	if err != nil {
		t.Fatalf("redirects: List no redirect found :%s", err)
	}

	if len(results) < 1 {
		t.Fatalf("redirects: List no redirects found :%s", err)
	}

}

// Test Update method
func TestUpdateRedirect(t *testing.T) {

	// Get the last redirect (created in TestCreateRedirect above)
	redirect, err := FindFirst("new_url=?", testURL)
	if err != nil {
		t.Fatalf("redirects: Update no redirect found :%s", err)
	}

	url := "/bar"
	redirectParams := map[string]string{"new_url": url}
	err = redirect.Update(redirectParams)
	if err != nil {
		t.Fatalf("redirects: Update redirect failed :%s", err)
	}

	// Fetch the redirect again from db
	redirect, err = Find(redirect.ID)
	if err != nil {
		t.Fatalf("redirects: Update redirect fetch failed :%s", redirect.NewURL)
	}

	if redirect.NewURL != url {
		t.Fatalf("redirects: Update redirect failed :%s", redirect.NewURL)
	}

}

// TestQuery tests trying to find published resources
func TestQuery(t *testing.T) {

	results, err := FindAll(Query())
	if err != nil {
		t.Fatalf("redirects: error getting redirects :%s", err)
	}
	if len(results) == 0 {
		t.Fatalf("redirects: published redirects not found :%s", err)
	}

	results, err = FindAll(Query().Where("id>=? AND id <=?", 0, 100))
	if err != nil || len(results) == 0 {
		t.Fatalf("redirects: no redirect found :%s", err)
	}
	if len(results) > 1 {
		t.Fatalf("redirects: more than one redirect found for where :%s", err)
	}

}

// Test Destroy method
func TestDestroyRedirect(t *testing.T) {

	results, err := FindAll(Query())
	if err != nil || len(results) == 0 {
		t.Fatalf("redirects: Destroy no redirect found :%s", err)
	}
	redirect := results[0]
	count := len(results)

	err = redirect.Destroy()
	if err != nil {
		t.Fatalf("redirects: Destroy redirect failed :%s", err)
	}

	// Check new length of redirects returned
	results, err = FindAll(Query())
	if err != nil {
		t.Fatalf("redirects: Destroy error getting results :%s", err)
	}

	// length should be one less than previous
	if len(results) != count-1 {
		t.Fatalf("redirects: Destroy redirect count wrong :%d", len(results))
	}

}

// TestAllowedParams should always return some params
func TestAllowedParams(t *testing.T) {
	if len(AllowedParams()) == 0 {
		t.Fatalf("redirects: no allowed params")
	}
}
