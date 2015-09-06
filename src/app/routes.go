package app

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"
)

// Define routes for this app
func setupRoutes(r *router.Router) {

	// Set the default file handler
	r.FileHandler = fileHandler

	// Add a files route to handle static images under files - nginx deals with this in production
	r.Add("/files/{path:.*}", fileHandler)

	// Add the home page route
	r.Add("/", HandleShowHome)

}

// Default static file handler, handles assets too
func fileHandler(context router.Context) {

	if serveAsset(context) {
		return
	}

	serveFile(context)

}

func notAuthorizedHandler(context router.Context) {
	view := view.New(context)
	view.RenderStatus(context, http.StatusUnauthorized)
}

func notFoundHandler(context router.Context) {
	view.RenderStatus(context, http.StatusNotFound)
}

// Default file handler, used in development - in production serve with nginx
func serveFile(context router.Context) {
	// Assuming we're running from the root of the website
	localPath := "./public" + path.Clean(context.Path())

	if _, err := os.Stat(localPath); err != nil {
		if os.IsNotExist(err) {
			notFoundHandler(context)
			return
		}

		// For other errors return not authorised
		notAuthorizedHandler(context)
		return
	}

	// If the file exists and we can access it, serve it
	http.ServeFile(context, context.Request(), localPath)
}

// Handle serving assets in dev (if we can) - return true on success
func serveAsset(context router.Context) bool {
	p := path.Clean(context.Path())
	// It must be under /assets, or we don't serve
	if !strings.HasPrefix(p, "/assets/") {
		return false
	}

	// Try to find an asset in our list
	f := appAssets.File(path.Base(p))
	if f == nil {
		return false
	}

	localPath := "./" + f.LocalPath()
	http.ServeFile(context, context.Request(), localPath)
	return true
}
