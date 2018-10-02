// Package parser defines an interface for parsers (creating templates) and templates (rendering content), and defines a base template type which conforms to both interfaces and can be included in any templates
package parser

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Scanner scans paths for templates and creates a representation of each using parsers
type Scanner struct {
	// A map of all templates keyed by path name
	Templates map[string]Template

	// A set of parsers (in order) with which to parse templates
	Parsers []Parser

	// A set of paths (in order) from which to load templates
	Paths []string

	// Helpers is a list of helper functions
	Helpers FuncMap

	// rootPath is used to store the root path during scans
	rootPath string
}

// NewScanner creates a new template scanner
func NewScanner(paths []string, helpers FuncMap) (*Scanner, error) {
	s := &Scanner{
		Helpers:   helpers,
		Paths:     paths,
		Templates: make(map[string]Template),
		Parsers:   []Parser{new(JSONTemplate), new(HTMLTemplate), new(TextTemplate)},
	}

	return s, nil
}

// ScanPath scans a path for template files, including sub-paths
func (s *Scanner) ScanPath(root string) error {

	// Store the rootPath - used in walkFunc
	s.rootPath = path.Clean(root)

	// Store current path, and change to root path
	// so that template includes use relative paths from root
	// this may not be necc. any more, test removing it
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Change dir to the rootPath so that paths are relative
	err = os.Chdir(s.rootPath)
	if err != nil {
		return err
	}

	err = filepath.Walk(".", s.walkFunc)
	if err != nil {
		return err
	}

	// Change back to original path
	err = os.Chdir(pwd)
	if err != nil {
		return err
	}

	return nil
}

// walkFunc handles files from filepath.Walk in ScanPath
// It follows symlinks where encountered by recursing
func (s *Scanner) walkFunc(path string, info os.FileInfo, err error) error {

	// If an error occurred, report it
	if err != nil {
		return err
	}

	// Check if this is a symlink, if so recurse
	// This assumes that the structure at linkedPath exactly mirrors that at path
	if isSymlink(info) {

		// Find the linked path
		linkedPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return fmt.Errorf("error reading symbolic link: %s", err)
		}

		// Calculate a new temp root path, based on linked path
		// trimmed of the original path
		// This assumes that the structure at linkedPath exactly mirrors that at path
		newRoot := strings.TrimSuffix(linkedPath, path)

		// Store current dir
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// Change dir to the linked path container
		err = os.Chdir(newRoot)
		if err != nil {
			return err
		}

		// filepath.Walk at the location, based on the newRoot
		err = filepath.Walk(path, s.walkFunc)

		// Change dir back to pwd
		err = os.Chdir(pwd)
		if err != nil {
			return err
		}

		return err
	}

	// Deal with files, directories we return nil error to recurse on them
	if !info.IsDir() {
		// Ask parsers in turn to handle the file - first one to claim it wins
		for _, p := range s.Parsers {
			if p.CanParseFile(path) {

				fullpath := filepath.Join(s.rootPath, path)
				t, err := p.NewTemplate(fullpath, path)
				if err != nil {
					return err
				}

				s.Templates[path] = t
				return nil
			}
		}

	}

	return nil
}

// isSymlink returns true if this is a symlink
func isSymlink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}

// ScanPaths resets template list and rescans all template paths
func (s *Scanner) ScanPaths() error {
	// Make sure templates is empty
	s.Templates = make(map[string]Template)

	// Set up the parsers
	for _, p := range s.Parsers {
		err := p.Setup(s.Helpers)
		if err != nil {
			return err
		}
	}

	// Scan paths again
	for _, p := range s.Paths {
		err := s.ScanPath(p)
		if err != nil {
			return err
		}
	}

	// Now parse and finalize templates
	for _, t := range s.Templates {
		err := t.Parse()
		if err != nil {
			return err
		}
	}

	// Now finalize templates
	for _, t := range s.Templates {
		err := t.Finalize(s.Templates)
		if err != nil {
			return err
		}
	}

	return nil
}

// PATH UTILITIES

// dotFile returns true if the file path supplied a dot file?
func dotFile(p string) bool {
	return strings.HasPrefix(path.Base(p), ".")
}

// suffix returns true if the path have this suffix (ignoring dotfiles)?
func suffix(p string, suffix string) bool {
	if dotFile(p) {
		return false
	}
	return strings.HasSuffix(p, suffix)
}

// suffixes returns true if the path has these suffixes (ignoring dotfiles)?
func suffixes(p string, suffixes []string) bool {
	if dotFile(p) {
		return false
	}

	for _, s := range suffixes {
		if strings.HasSuffix(p, s) {
			return true
		}
	}

	return false
}
