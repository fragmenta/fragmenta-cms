package pageactions

import (
	"fmt"
	"io/ioutil"

	"strings"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/pages"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

// HandleShowSetup shows the setup page at /fragmenta/setup
func HandleShowSetup(context router.Context) error {

	// Setup context for template
	view := view.New(context)

	// If we have pages or users already, do not proceed
	if !missingUsersAndPages() {
		return router.NotAuthorizedError(nil)
	}

	view.Template("pages/views/setup.html.got")
	return view.Render()
}

// HandleSetup responds to a POST at /fragmenta/setup
// by creating our first user and page
func HandleSetup(context router.Context) error {

	// If we have pages or users already, do not proceed
	if !missingUsersAndPages() {
		return router.NotAuthorizedError(nil)
	}

	// Take the details given and create the first user
	params := map[string]string{
		"email":    context.Param("email"),
		"password": context.Param("password"),
		"name":     nameFromEmail(context.Param("email")),
		"status":   "100",
		"role":     "100",
		"title":    "Administrator",
	}

	uid, err := users.Create(params)
	if err != nil {
		return router.InternalError(err)
	}
	context.Logf("#info Created user #%d", uid)
	user, err := users.Find(uid)
	if err != nil {
		return router.InternalError(err)
	}
	// Login this user automatically - save cookie
	session, err := auth.Session(context, context.Request())
	if err != nil {
		return router.InternalError(err)
	}
	context.Logf("#info Automatic login for first user: %d %s", user.Id, user.Email)
	session.Set(auth.SessionUserKey, fmt.Sprintf("%d", user.Id))
	session.Save(context)

	// Load our welcomepage template html
	// and put it into the text field of a new page with id 1

	welcomeText, err := ioutil.ReadFile("src/pages/views/welcome.html.got")
	if err != nil {
		return router.InternalError(err)
	}

	params = map[string]string{
		"status": "100",
		"name":   "Fragmenta",
		"url":    "/",
		"text":   string(welcomeText),
	}
	_, err = pages.Create(params)
	if err != nil {
		return router.InternalError(err)
	}

	// Create another couple of simple pages as examples (about and privacy)
	params = map[string]string{
		"status": "100",
		"name":   "About Us",
		"url":    "/about",
		"text":   "<section class=\"narrow\"><h1>About us</h1><p>About us</p></section>",
	}
	_, err = pages.Create(params)
	if err != nil {
		return router.InternalError(err)
	}
	params = map[string]string{
		"status": "100",
		"name":   "Privacy Policy",
		"url":    "/privacy",
		"text":   "<section class=\"narrow\"><h1>Privacy Policy</h1><p>We respect your privacy.</p></section>",
	}
	_, err = pages.Create(params)
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect back to the newly populated home page
	return router.Redirect(context, "/")
}

// nameFromEmail grabs a name from an email address
func nameFromEmail(e string) string {
	// Split email on @, and separate by removing . or _
	parts := strings.Split(e, "@")
	if len(parts) > 0 {
		n := strings.Replace(parts[0], ".", " ", -1)
		n = strings.Replace(n, "_", " ", -1)
		return n
	}

	return e
}

// missingUsersAndPages returns true if we have 0 users and 0 pages
func missingUsersAndPages() bool {

	pageCount, err := pages.Query().Count()
	if err != nil {
		return true
	}

	userCount, err := users.Query().Count()
	if err != nil {
		return true
	}

	if pageCount > 0 || userCount > 0 {
		return false
	}

	return true
}
