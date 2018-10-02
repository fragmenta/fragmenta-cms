// Package parser defines an interface for parsers which the base template conforms to
package parser

// FuncMap is a map of functions
type FuncMap map[string]interface{}

// Parser loads template files, and returns a template suitable for rendering content
type Parser interface {
	// Setup is called once on setup of a parser
	Setup(helpers FuncMap) error

	// Can this parser handle this file?
	CanParseFile(path string) bool

	// Parse the file given and return a compiled template
	NewTemplate(fullpath, path string) (Template, error)
}
