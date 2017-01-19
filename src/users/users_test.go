// Tests for the users package
package users

import (
	"testing"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
)

var testUserName = "'fué ';'\""

func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(2)
	if err != nil {
		t.Errorf("users: Setup db failed %s", err)
	}
}

// Test Create method
func TestCreateUser(t *testing.T) {
	userParams := map[string]string{
		"name":   testUserName,
		"status": "100",
	}
	id, err := New().Create(userParams)
	if err != nil {
		t.Errorf("users: Create user failed :%s", err)
	}

	user, err := Find(id)
	if err != nil {
		t.Errorf("users: Create user find failed")
	}

	if user.Name != testUserName {
		t.Errorf("users: Create user name failed expected:%s got:%s", testUserName, user.Name)
	}

}

// Test Index (List) method
func TestListUsers(t *testing.T) {

	// Get all users (we should have at least one)
	results, err := FindAll(Query())
	if err != nil {
		t.Errorf("users: List no user found :%s", err)
	}

	if len(results) < 1 {
		t.Errorf("users: List no users found :%s", err)
	}

}

// Test Update method
func TestUpdateUser(t *testing.T) {

	// Get the last user (created in TestCreateUser above)
	user, err := FindFirst("name=?", testUserName)
	if err != nil {
		t.Errorf("users: Update no user found :%s", err)
	}

	name := "bar"
	userParams := map[string]string{"name": name}
	err = user.Update(userParams)
	if err != nil {
		t.Errorf("users: Update user failed :%s", err)
	}

	// Fetch the user again from db
	user, err = Find(user.ID)
	if err != nil {
		t.Errorf("users: Update user fetch failed :%s", user.Name)
	}

	if user.Name != name {
		t.Errorf("users: Update user failed :%s", user.Name)
	}

}

// Test Destroy method
func TestDestroyUser(t *testing.T) {

	results, err := FindAll(Query())
	if err != nil || len(results) == 0 {
		t.Errorf("users: Destroy no user found :%s", err)
	}
	user := results[0]
	count := len(results)

	err = user.Destroy()
	if err != nil {
		t.Errorf("users: Destroy user failed :%s", err)
	}

	// Check new length of users returned
	results, err = FindAll(Query())
	if err != nil {
		t.Errorf("users: Destroy error getting results :%s", err)
	}

	// length should be one less than previous
	if len(results) != count-1 {
		t.Errorf("users: Destroy user count wrong :%d", len(results))
	}

}

func TestQuery(t *testing.T) {

	results, err := FindAll(Published())
	if err != nil {
		t.Errorf("users: error getting users :%s", err)
	}
	if len(results) == 0 {
		t.Errorf("users: published users not found :%s", err)
	}

	results, err = FindAll(Query().Where("id>=? AND id <=?", 0, 100))
	if err != nil {
		t.Errorf("users: no user found :%s", err)
	}
	if len(results) > 2 {
		t.Errorf("users: more than 2 users found for where :%v", results)
	}

}

func TestRoles(t *testing.T) {
	u := MockAdmin()
	options := u.RoleOptions()
	if len(options) == 0 {
		t.Errorf("users: error creating role options")
	}

	if u.RoleDisplay() != "Administrator" {
		t.Errorf("users: error displaying role")
	}

	if !u.Admin() || u.Reader() || u.Anon() {
		t.Errorf("users: error testing role")
	}

	u = MockAnon()
	if u.Admin() || u.RoleID() > Anon || u.UserID() > 0 {
		t.Errorf("users: error testing anon role")
	}

}

// TestAllowedParams should always return some params
func TestAllowedParams(t *testing.T) {
	if len(AllowedParams()) == 0 {
		t.Errorf("users: no allowed params")
	}
}
