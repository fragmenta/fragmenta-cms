// Tests for the tags package
package tags

import (
	"testing"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
)

var testName = "foo"

func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(2)
	if err != nil {
		t.Errorf("tags: Setup db failed %s", err)
	}
}

// Test Create method
func TestCreateTags(t *testing.T) {
	tagParams := map[string]string{
		"name":   testName,
		"status": "100",
	}

	id, err := New().Create(tagParams)
	if err != nil {
		t.Errorf("tags: Create tag failed :%s", err)
	}

	tag, err := Find(id)
	if err != nil {
		t.Errorf("tags: Create tag find failed")
	}

	if tag.Name != testName {
		t.Errorf("tags: Create tag name failed expected:%s got:%s", testName, tag.Name)
	}

}

// Test Index (List) method
func TestListTags(t *testing.T) {

	// Get all tags (we should have at least one)
	results, err := FindAll(Query())
	if err != nil {
		t.Errorf("tags: List no tag found :%s", err)
	}

	if len(results) < 1 {
		t.Errorf("tags: List no tags found :%s", err)
	}

}

// Test Update method
func TestUpdateTags(t *testing.T) {

	// Get the last tag (created in TestCreateTags above)
	tag, err := FindFirst("name=?", testName)
	if err != nil {
		t.Errorf("tags: Update no tag found :%s", err)
	}

	name := "bar"
	tagParams := map[string]string{"name": name}
	err = tag.Update(tagParams)
	if err != nil {
		t.Errorf("tags: Update tag failed :%s", err)
	}

	// Fetch the tag again from db
	tag, err = Find(tag.ID)
	if err != nil {
		t.Errorf("tags: Update tag fetch failed :%s", tag.Name)
	}

	if tag.Name != name {
		t.Errorf("tags: Update tag failed :%s", tag.Name)
	}

}

// TestQuery tests trying to find published resources
func TestQuery(t *testing.T) {

	results, err := FindAll(Published())
	if err != nil {
		t.Errorf("tags: error getting tags :%s", err)
	}
	if len(results) == 0 {
		t.Errorf("tags: published tags not found :%s", err)
	}

	results, err = FindAll(Query().Where("id>=? AND id <=?", 0, 100))
	if err != nil || len(results) == 0 {
		t.Errorf("tags: no tag found :%s", err)
	}
	if len(results) > 1 {
		t.Errorf("tags: more than one tag found for where :%s", err)
	}

}

// Test Destroy method
func TestDestroyTags(t *testing.T) {

	results, err := FindAll(Query())
	if err != nil || len(results) == 0 {
		t.Errorf("tags: Destroy no tag found :%s", err)
	}
	tag := results[0]
	count := len(results)

	err = tag.Destroy()
	if err != nil {
		t.Errorf("tags: Destroy tag failed :%s", err)
	}

	// Check new length of tags returned
	results, err = FindAll(Query())
	if err != nil {
		t.Errorf("tags: Destroy error getting results :%s", err)
	}

	// length should be one less than previous
	if len(results) != count-1 {
		t.Errorf("tags: Destroy tag count wrong :%d", len(results))
	}

}

// TestAllowedParams should always return some params
func TestAllowedParams(t *testing.T) {
	if len(AllowedParams()) == 0 {
		t.Errorf("tags: no allowed params")
	}
}
