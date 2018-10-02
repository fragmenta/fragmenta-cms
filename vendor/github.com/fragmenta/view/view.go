// Package view provides methods for rendering templates, and helper functions for golang views
package view

import (
	"fmt"
	"sync"

	"github.com/fragmenta/view/helpers"
	"github.com/fragmenta/view/parser"
)

// Production is true if this server is running in production mode
var Production bool

// The scanner is a private type used for scanning templates
var scanner *parser.Scanner

// This mutex guards the pkg scanner variable during reload and access
// it is only neccessary because of hot reload during development
var mu sync.RWMutex

// Helpers is a list of functions available in templates
var Helpers parser.FuncMap

func init() {
	Helpers = DefaultHelpers()
}

// LoadTemplates loads our templates from ./src, and assigns them to the package variable Templates
// This function is deprecated and will be removed, use LoadTemplatesAtPaths instead
func LoadTemplates() error {
	return LoadTemplatesAtPaths([]string{"src"}, Helpers)
}

// DefaultHelpers returns a default set of helpers for the app,
// which can then be extended/replaced. Helper functions may not be changed
// after LoadTemplates is called, as reloading is required if they change.
func DefaultHelpers() parser.FuncMap {
	funcs := make(parser.FuncMap)

	// HEAD helpers
	funcs["style"] = helpers.Style
	funcs["script"] = helpers.Script
	funcs["dev"] = func() bool { return !Production }

	// HTML helpers
	funcs["html"] = helpers.HTML
	funcs["htmlattr"] = helpers.HTMLAttribute
	funcs["url"] = helpers.URL

	funcs["sanitize"] = helpers.Sanitize
	funcs["strip"] = helpers.Strip
	funcs["truncate"] = helpers.Truncate

	// XML helpers
	funcs["xmlpreamble"] = helpers.XMLPreamble

	// Form helpers
	funcs["field"] = helpers.Field
	funcs["datefield"] = helpers.DateField
	funcs["textarea"] = helpers.TextArea
	funcs["select"] = helpers.Select
	funcs["selectarray"] = helpers.SelectArray
	funcs["optionsforselect"] = helpers.OptionsForSelect

	funcs["utcdate"] = helpers.UTCDate
	funcs["utctime"] = helpers.UTCTime
	funcs["utcnow"] = helpers.UTCNow
	funcs["date"] = helpers.Date
	funcs["time"] = helpers.Time
	funcs["numberoptions"] = helpers.NumberOptions

	// CSV helper
	funcs["csv"] = helpers.CSV

	// String helpers
	funcs["blank"] = helpers.Blank
	funcs["exists"] = helpers.Exists

	// Math helpers
	funcs["mod"] = helpers.Mod
	funcs["odd"] = helpers.Odd
	funcs["add"] = helpers.Add
	funcs["subtract"] = helpers.Subtract

	// Array functions
	funcs["array"] = helpers.Array
	funcs["append"] = helpers.Append
	funcs["contains"] = helpers.Contains

	// Map functions
	funcs["map"] = helpers.Map
	funcs["set"] = helpers.Set
	funcs["setif"] = helpers.SetIf
	funcs["empty"] = helpers.Empty

	// Numeric helpers - clean up and accept currency and other options in centstoprice
	funcs["centstobase"] = helpers.CentsToBase
	funcs["centstoprice"] = helpers.CentsToPrice
	funcs["centstopriceshort"] = helpers.CentsToPriceShort
	funcs["pricetocents"] = helpers.PriceToCents

	return funcs
}

// LoadTemplatesAtPaths loads our templates given the paths provided
func LoadTemplatesAtPaths(paths []string, helpers parser.FuncMap) error {

	mu.Lock()
	defer mu.Unlock()

	// Scan all templates within the given paths, using the helpers provided
	var err error
	scanner, err = parser.NewScanner(paths, helpers)
	if err != nil {
		return err
	}

	err = scanner.ScanPaths()
	if err != nil {
		return err
	}

	return nil
}

// ReloadTemplates reloads the templates for our scanner
func ReloadTemplates() error {
	mu.Lock()
	defer mu.Unlock()
	return scanner.ScanPaths()
}

// PrintTemplates prints out our list of templates for debug
func PrintTemplates() {
	mu.RLock()
	defer mu.RUnlock()
	for k := range scanner.Templates {
		fmt.Printf("%s\n", k)
	}
	fmt.Printf("Finished scan of templates\n")
}
