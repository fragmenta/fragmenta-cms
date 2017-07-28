// Tests for the redirects package
package redirects

import (
	"testing"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
)

var testName = "foo"

func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(2)
	if err != nil {
		t.Fatalf("redirects: Setup db failed %s", err)
	}
}

// Test Create method
func TestCreateRedirect(t *testing.T) {
	redirectParams := map[string]string{
		"name":   testName,
		"status": "100",
	}

	id, err := New().Create(redirectParams)
	if err != nil {
		t.Fatalf("redirects: Create redirect failed :%s", err)
	}

	redirect, err := Find(id)
	if err != nil {
		t.Fatalf("redirects: Create redirect find failed")
	}

	if redirect.Name != testName {
		t.Fatalf("redirects: Create redirect name failed expected:%s got:%s", testName, redirect.Name)
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
	redirect, err := FindFirst("name=?", testName)
	if err != nil {
		t.Fatalf("redirects: Update no redirect found :%s", err)
	}

	name := "bar"
	redirectParams := map[string]string{"name": name}
	err = redirect.Update(redirectParams)
	if err != nil {
		t.Fatalf("redirects: Update redirect failed :%s", err)
	}

	// Fetch the redirect again from db
	redirect, err = Find(redirect.ID)
	if err != nil {
		t.Fatalf("redirects: Update redirect fetch failed :%s", redirect.Name)
	}

	if redirect.Name != name {
		t.Fatalf("redirects: Update redirect failed :%s", redirect.Name)
	}

}

// TestQuery tests trying to find published resources
func TestQuery(t *testing.T) {

	results, err := FindAll(Published())
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
