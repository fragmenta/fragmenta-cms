package authorise

import (
	"strconv"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/router"

	"github.com/fragmenta/fragmenta-cms/src/users"
)

// CurrentUser returns the saved user (or an empty anon user) for the current session cookie
// Strictly speaking this should be authenticate.User
func CurrentUser(context router.Context) *users.User {

	// First check if the user has already been set on context, if so return it
	if context.Get("current_user") != nil {
		return context.Get("current_user").(*users.User)
	}

	// Start with an anon user by default (role 0, id 0)
	user := &users.User{}

	// Build the session from the secure cookie, or create a new one
	session, err := auth.Session(context, context.Request())
	if err != nil {
		context.Logf("#error problem retrieving session")
		return user
	}

	// Fetch the current user record if we have one recorded in the session
	var id int64
	ids := session.Get(auth.SessionUserKey)
	if len(ids) > 0 {
		id, err = strconv.ParseInt(ids, 10, 64)
		if err != nil {
			context.Logf("#error Error decoding session user key:%s\n", err)
			return user
		}
	}

	if id != 0 {
		u, err := users.Find(id)
		if err != nil {
			context.Logf("#info User not found from session id:%d\n", id)
			return user
		}
		user = u
	}

	return user
}
