package server

import (
	"fmt"
	"net/http"
	"time"
)

//

// daysToSeconds is a time constant for converting days to seconds
const daysToSeconds = 86400

// Date format is the preferred date format for the Expires header
const dateFormat = "Mon, 2 Jan 2006 15:04:05 MST"

// AddCacheHeaders adds Cache-Control, Expires and Etag headers
// using the age in days and content hash provided
func AddCacheHeaders(w http.ResponseWriter, days int, hash string) {
	// Cache for the given age in days
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age:%d", days*daysToSeconds))

	// Set an expires header of form Mon Jan 2 15:04:05 -0700 MST 2006
	w.Header().Set("Expires", time.Now().AddDate(0, 0, days).UTC().Format(dateFormat))

	// For etag send the hash given
	w.Header().Set("ETag", fmt.Sprintf("\"%s\"", hash))
}
