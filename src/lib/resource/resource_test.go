package resource

import (
	"testing"
)

var r = Base{ID: 99, TableName: "images", KeyName: "id"}

func TestValidate(t *testing.T) {

	// TEST NIL VALUES in db
	// most validation is to ensure nil db columns return zero values rather than panic

	// Test against nil bool
	if ValidateBoolean(nil) != false {
		t.Fatalf("Validate does not match expected:%v got:%v", false, ValidateBoolean(nil))
	}

	// Test against nil float
	if ValidateFloat(nil) != 0 {
		t.Fatalf("Validate does not match expected:%v got:%v", 0, ValidateFloat(nil))
	}

	// Test against nil int64
	if ValidateInt(nil) != 0 {
		t.Fatalf("Validate does not match expected:%v got:%v", 0, ValidateInt(nil))
	}

	// Test against nil time
	if !ValidateTime(nil).IsZero() {
		t.Fatalf("Validate does not match expected:%v got:%v", "zero time", ValidateTime(nil))
	}

	// Test against nil string
	if ValidateString(nil) != "" {
		t.Fatalf("Validate does not match expected:%v got:%v", "", ValidateString(nil))
	}

	// TEST VALUES

	if ValidateBoolean(true) != true { // yes, I know!
		t.Fatalf("Validate does not match expected:%v got:%v", true, ValidateBoolean(true))
	}

	// Test against range of ints as sanity check on casts
	ints := []int{99, 5, 0, 1, 1110011, -1200}
	for _, i := range ints {
		if ValidateInt(i) != int64(i) {
			t.Fatalf("Validate float does not match expected:%v got:%v", i, ValidateInt(i))
		}
	}

	// Test against range of floats as sanity check
	floats := []float64{5.0, 0.0, 0.0001, -0.40}
	for _, f := range floats {
		if ValidateFloat(f) != f {
			t.Fatalf("Validate float does not match expected:%v got:%v", f, ValidateFloat(f))
		}
	}
	// Check success of cast of int column values to floats
	for _, f := range ints {
		if ValidateFloat(f) != float64(f) {
			t.Fatalf("Validate float does not match expected:%v got:%v", f, ValidateFloat(f))
		}
	}

}

// TestValidateParams tests we remove params correctly when not authorised
// the default set is an empty set
func TestValidateParams(t *testing.T) {
	allowed := []string{"name"}
	params := map[string]string{"name_asdfasdf_name": "asdf", "name": "foo", "bar": "baz"}
	params = r.ValidateParams(params, allowed)
	if len(params) > 1 {
		t.Fatalf("Validate params does not match expected:%v got:%v", "[]", params)
	}
}

// We should probably have an in-memory adapter for query which lets us test creation etc easily

// TestURLs tests the url functions in urls.go
func TestURLs(t *testing.T) {

	expected := "/images"
	if r.IndexURL() != expected {
		t.Fatalf("URL does not match expected:%s got:%s", expected, r.IndexURL())
	}
	expected = "/images/create"
	if r.CreateURL() != expected {
		t.Fatalf("URL does not match expected:%s got:%s", expected, r.CreateURL())
	}
	expected = "/images/99/update"
	if r.UpdateURL() != expected {
		t.Fatalf("URL does not match expected:%s got:%s", expected, r.UpdateURL())
	}
	expected = "/images/99"
	if r.ShowURL() != expected {
		t.Fatalf("URL does not match expected:%s got:%s", expected, r.ShowURL())
	}
	expected = "/images/99/destroy"
	if r.DestroyURL() != expected {
		t.Fatalf("URL does not match expected:%s got:%s", expected, r.DestroyURL())
	}

}
