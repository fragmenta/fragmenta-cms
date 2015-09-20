package app

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/fragmenta/router"

	"github.com/fragmenta/view"
)

// Default static file handler, handles assets too
func fileHandler(context router.Context) error {

	err := serveAsset(context)
	if err == nil {
		return nil // return on success only for assets
	}

	// Finally try serving a file from public
	return serveFile(context)
}

// Default file handler, used in development - in production serve with nginx
func serveFile(context router.Context) error {
	// Assuming we're running from the root of the website
	localPath := "./public" + path.Clean(context.Path())

	if _, err := os.Stat(localPath); err != nil {
		// If file not found return error
		if os.IsNotExist(err) {
			return router.NotFoundError(err)
		}

		// For other errors return not authorised
		return router.NotAuthorizedError(err)
	}

	// If the file exists and we can access it, serve it
	http.ServeFile(context, context.Request(), localPath)
	return nil
}

// Handle serving assets in dev (if we can) - return true on success
func serveAsset(context router.Context) error {
	p := path.Clean(context.Path())

	// It must be under /assets, or we don't serve
	if !strings.HasPrefix(p, "/assets/") {
		return router.NotFoundError(nil)
	}

	// Try to find an asset in our list
	f := appAssets.File(path.Base(p))
	if f == nil {
		return router.NotFoundError(nil)
	}

	localPath := "./" + f.LocalPath()
	http.ServeFile(context, context.Request(), localPath)
	return nil
}

// errHandler renders an error using error templates if available
func errHandler(context router.Context, e error) {

	// Cast the error to a status error if it is one, if not wrap it in a Status 500 error
	err := router.ToStatusError(e)

	view := view.New(context)

	view.AddKey("title", err.Title)
	view.AddKey("message", err.Message)

	if !context.Production() {
		view.AddKey("status", err.Status)
		view.AddKey("file", err.FileLine())
		view.AddKey("error", err.Err)
	}

	view.Template("app/views/error.html.got")

	context.Logf("#error %s\n", err)
	view.Render()
}
