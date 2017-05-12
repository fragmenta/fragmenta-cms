// Tests for the pages package
package pages

import (
	"testing"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
)

var testName = "foo"

func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(2)
	if err != nil {
		t.Fatalf("pages: Setup db failed %s", err)
	}
}

// Test Create method
func TestCreatePage(t *testing.T) {
	pageParams := map[string]string{
		"name":   testName,
		"status": "100",
	}

	id, err := New().Create(pageParams)
	if err != nil {
		t.Fatalf("pages: Create page failed :%s", err)
	}

	page, err := Find(id)
	if err != nil {
		t.Fatalf("pages: Create page find failed")
	}

	if page.Name != testName {
		t.Fatalf("pages: Create page name failed expected:%s got:%s", testName, page.Name)
	}

}

// Test Index (List) method
func TestListPage(t *testing.T) {

	// Get all pages (we should have at least one)
	results, err := FindAll(Query())
	if err != nil {
		t.Fatalf("pages: List no page found :%s", err)
	}

	if len(results) < 1 {
		t.Fatalf("pages: List no pages found :%s", err)
	}

}

// Test Update method
func TestUpdatePage(t *testing.T) {

	// Get the last page (created in TestCreatePage above)
	page, err := FindFirst("name=?", testName)
	if err != nil {
		t.Fatalf("pages: Update no page found :%s", err)
	}

	name := "bar"
	pageParams := map[string]string{"name": name}
	err = page.Update(pageParams)
	if err != nil {
		t.Fatalf("pages: Update page failed :%s", err)
	}

	// Fetch the page again from db
	page, err = Find(page.ID)
	if err != nil {
		t.Fatalf("pages: Update page fetch failed :%s", page.Name)
	}

	if page.Name != name {
		t.Fatalf("pages: Update page failed :%s", page.Name)
	}

}

// TestQuery tests trying to find published resources
func TestQuery(t *testing.T) {

	results, err := FindAll(Published())
	if err != nil {
		t.Fatalf("pages: error getting pages :%s", err)
	}
	if len(results) == 0 {
		t.Fatalf("pages: published pages not found :%s", err)
	}

	results, err = FindAll(Query().Where("id>=? AND id <=?", 0, 100))
	if err != nil || len(results) == 0 {
		t.Fatalf("pages: no page found :%s", err)
	}
	if len(results) > 1 {
		t.Fatalf("pages: more than one page found for where :%s", err)
	}

}

// Test Destroy method
func TestDestroyPage(t *testing.T) {

	results, err := FindAll(Query())
	if err != nil || len(results) == 0 {
		t.Fatalf("pages: Destroy no page found :%s", err)
	}
	page := results[0]
	count := len(results)

	err = page.Destroy()
	if err != nil {
		t.Fatalf("pages: Destroy page failed :%s", err)
	}

	// Check new length of pages returned
	results, err = FindAll(Query())
	if err != nil {
		t.Fatalf("pages: Destroy error getting results :%s", err)
	}

	// length should be one less than previous
	if len(results) != count-1 {
		t.Fatalf("pages: Destroy page count wrong :%d", len(results))
	}

}

// TestAllowedParams should always return some params
func TestAllowedParams(t *testing.T) {
	if len(AllowedParams()) == 0 {
		t.Fatalf("pages: no allowed params")
	}
}
