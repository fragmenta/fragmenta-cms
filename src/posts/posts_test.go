// Tests for the posts package
package posts

import (
	"testing"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
)

var testName = "foo"

func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(2)
	if err != nil {
		t.Fatalf("posts: Setup db failed %s", err)
	}
}

// Test Create method
func TestCreatePost(t *testing.T) {
	postParams := map[string]string{
		"name":   testName,
		"status": "100",
	}

	id, err := New().Create(postParams)
	if err != nil {
		t.Fatalf("posts: Create post failed :%s", err)
	}

	post, err := Find(id)
	if err != nil {
		t.Fatalf("posts: Create post find failed")
	}

	if post.Name != testName {
		t.Fatalf("posts: Create post name failed expected:%s got:%s", testName, post.Name)
	}

}

// Test Index (List) method
func TestListPost(t *testing.T) {

	// Get all posts (we should have at least one)
	results, err := FindAll(Query())
	if err != nil {
		t.Fatalf("posts: List no post found :%s", err)
	}

	if len(results) < 1 {
		t.Fatalf("posts: List no posts found :%s", err)
	}

}

// Test Update method
func TestUpdatePost(t *testing.T) {

	// Get the last post (created in TestCreatePost above)
	post, err := FindFirst("name=?", testName)
	if err != nil {
		t.Fatalf("posts: Update no post found :%s", err)
	}

	name := "bar"
	postParams := map[string]string{"name": name}
	err = post.Update(postParams)
	if err != nil {
		t.Fatalf("posts: Update post failed :%s", err)
	}

	// Fetch the post again from db
	post, err = Find(post.ID)
	if err != nil {
		t.Fatalf("posts: Update post fetch failed :%s", post.Name)
	}

	if post.Name != name {
		t.Fatalf("posts: Update post failed :%s", post.Name)
	}

}

// TestQuery tests trying to find published resources
func TestQuery(t *testing.T) {

	results, err := FindAll(Published())
	if err != nil {
		t.Fatalf("posts: error getting posts :%s", err)
	}
	if len(results) == 0 {
		t.Fatalf("posts: published posts not found :%s", err)
	}

	results, err = FindAll(Query().Where("id>=? AND id <=?", 0, 100))
	if err != nil || len(results) == 0 {
		t.Fatalf("posts: no post found :%s", err)
	}
	if len(results) > 1 {
		t.Fatalf("posts: more than one post found for where :%s", err)
	}

}

// Test Destroy method
func TestDestroyPost(t *testing.T) {

	results, err := FindAll(Query())
	if err != nil || len(results) == 0 {
		t.Fatalf("posts: Destroy no post found :%s", err)
	}
	post := results[0]
	count := len(results)

	err = post.Destroy()
	if err != nil {
		t.Fatalf("posts: Destroy post failed :%s", err)
	}

	// Check new length of posts returned
	results, err = FindAll(Query())
	if err != nil {
		t.Fatalf("posts: Destroy error getting results :%s", err)
	}

	// length should be one less than previous
	if len(results) != count-1 {
		t.Fatalf("posts: Destroy post count wrong :%d", len(results))
	}

}

// TestAllowedParams should always return some params
func TestAllowedParams(t *testing.T) {
	if len(AllowedParams()) == 0 {
		t.Fatalf("posts: no allowed params")
	}
}
