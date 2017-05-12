// Tests for the images package
package images

import (
	"testing"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
)

var testName = "foo"

func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(2)
	if err != nil {
		t.Fatalf("images: Setup db failed %s", err)
	}
}

// Test Create method
func TestCreateImage(t *testing.T) {
	imageParams := map[string]string{
		"name":   testName,
		"status": "100",
	}

	id, err := New().Create(imageParams)
	if err != nil {
		t.Fatalf("images: Create image failed :%s", err)
	}

	image, err := Find(id)
	if err != nil {
		t.Fatalf("images: Create image find failed")
	}

	if image.Name != testName {
		t.Fatalf("images: Create image name failed expected:%s got:%s", testName, image.Name)
	}

}

// Test Index (List) method
func TestListImage(t *testing.T) {

	// Get all images (we should have at least one)
	results, err := FindAll(Query())
	if err != nil {
		t.Fatalf("images: List no image found :%s", err)
	}

	if len(results) < 1 {
		t.Fatalf("images: List no images found :%s", err)
	}

}

// Test Update method
func TestUpdateImage(t *testing.T) {

	// Get the last image (created in TestCreateImage above)
	image, err := FindFirst("name=?", testName)
	if err != nil {
		t.Fatalf("images: Update no image found :%s", err)
	}

	name := "bar"
	imageParams := map[string]string{"name": name}
	err = image.Update(imageParams)
	if err != nil {
		t.Fatalf("images: Update image failed :%s", err)
	}

	// Fetch the image again from db
	image, err = Find(image.ID)
	if err != nil {
		t.Fatalf("images: Update image fetch failed :%s", image.Name)
	}

	if image.Name != name {
		t.Fatalf("images: Update image failed :%s", image.Name)
	}

}

// TestQuery tests trying to find published resources
func TestQuery(t *testing.T) {

	results, err := FindAll(Published())
	if err != nil {
		t.Fatalf("images: error getting images :%s", err)
	}
	if len(results) == 0 {
		t.Fatalf("images: published images not found :%s", err)
	}

	results, err = FindAll(Query().Where("id>=? AND id <=?", 0, 100))
	if err != nil || len(results) == 0 {
		t.Fatalf("images: no image found :%s", err)
	}
	if len(results) > 1 {
		t.Fatalf("images: more than one image found for where :%s", err)
	}

}

// Test Destroy method
func TestDestroyImage(t *testing.T) {

	results, err := FindAll(Query())
	if err != nil || len(results) == 0 {
		t.Fatalf("images: Destroy no image found :%s", err)
	}
	image := results[0]
	count := len(results)

	err = image.Destroy()
	if err != nil {
		t.Fatalf("images: Destroy image failed :%s", err)
	}

	// Check new length of images returned
	results, err = FindAll(Query())
	if err != nil {
		t.Fatalf("images: Destroy error getting results :%s", err)
	}

	// length should be one less than previous
	if len(results) != count-1 {
		t.Fatalf("images: Destroy image count wrong :%d", len(results))
	}

}

// TestAllowedParams should always return some params
func TestAllowedParams(t *testing.T) {
	if len(AllowedParams()) == 0 {
		t.Fatalf("images: no allowed params")
	}
}
