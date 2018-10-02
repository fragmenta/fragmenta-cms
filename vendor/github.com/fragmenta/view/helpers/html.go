package helpers

import (
	"fmt"
	"strings"

	got "html/template"

	"github.com/kennygrant/sanitize"
)

// Style inserts a css tag
func Style(name string) got.HTML {
	return got.HTML(fmt.Sprintf("<link href=\"/assets/styles/%s.css\" media=\"all\" rel=\"stylesheet\" type=\"text/css\" />", EscapeURL(name)))
}

// Script inserts a script tag
func Script(name string) got.HTML {
	return got.HTML(fmt.Sprintf("<script src=\"/assets/scripts/%s.js\" type=\"text/javascript\"></script>", EscapeURL(name)))
}

// Escape escapes HTML using HTMLEscapeString
func Escape(s string) string {
	return got.HTMLEscapeString(s)
}

// EscapeURL escapes URLs using HTMLEscapeString
func EscapeURL(s string) string {
	return got.URLQueryEscaper(s)
}

// Link returns got.HTML with an anchor link given text and URL required
// Attributes (if supplied) should not contain user input
func Link(t string, u string, a ...string) got.HTML {
	attributes := ""
	if len(a) > 0 {
		attributes = strings.Join(a, " ")
	}
	return got.HTML(fmt.Sprintf("<a href=\"%s\" %s>%s</a>", Escape(u), Escape(attributes), Escape(t)))
}

// HTML returns a string (which must not contain user input) as go template HTML
func HTML(s string) got.HTML {
	return got.HTML(s)
}

// HTMLAttribute returns a string (which must not contain user input) as go template HTMLAttr
func HTMLAttribute(s string) got.HTMLAttr {
	return got.HTMLAttr(s)
}

// URL returns returns a string (which must not contain user input) as go template URL
func URL(s string) got.URL {
	return got.URL(s)
}

// Strip all html tags and returns as go template HTML
func Strip(s string) got.HTML {
	return got.HTML(sanitize.HTML(s))
}

// Sanitize the html, leaving only tags we consider safe (see the sanitize package for details and tests)
func Sanitize(s string) got.HTML {
	s, err := sanitize.HTMLAllowing(s)
	if err != nil {
		fmt.Printf("#error sanitizing html:%s", err)
		return got.HTML("")
	}
	return got.HTML(s)
}

// XMLPreamble returns an XML preamble as got.HTML,
// primarily to work around a bug in html/template which escapes <?
// see https://github.com/golang/go/issues/12496
func XMLPreamble() got.HTML {
	return `<?xml version="1.0" encoding="UTF-8"?>`
}
