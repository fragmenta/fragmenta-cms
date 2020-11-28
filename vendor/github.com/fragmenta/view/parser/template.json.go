package parser

import (
	"fmt"
	got "html/template"
	"io"
	"sync"
)

var (
	jsonMu          sync.RWMutex  // Shared mutex to go with shared template set, because of dev reloads
	jsonTemplateSet *got.Template // This is a shared template set for json templates
)

// JSONTemplate represents a template using go HTML/template
type JSONTemplate struct {
	BaseTemplate
}

// Setup performs one-time setup before parsing templates
func (t *JSONTemplate) Setup(helpers FuncMap) error {
	mu.Lock()
	defer mu.Unlock()
	jsonTemplateSet = got.New("").Funcs(got.FuncMap(helpers))
	return nil
}

// CanParseFile returns true if this template can parse this file
func (t *JSONTemplate) CanParseFile(path string) bool {
	allowed := []string{".json.got"}
	return suffixes(path, allowed)
}

// NewTemplate returns a new JSONTemplate
func (t *JSONTemplate) NewTemplate(fullpath, path string) (Template, error) {
	template := new(JSONTemplate)
	template.fullpath = fullpath
	template.path = path
	return template, nil
}

// Parse the template
func (t *JSONTemplate) Parse() error {
	mu.Lock()
	defer mu.Unlock()
	err := t.BaseTemplate.Parse()

	// Add to our template set
	if jsonTemplateSet.Lookup(t.Path()) == nil {
		_, err = jsonTemplateSet.New(t.path).Parse(t.Source())
	} else {
		err = fmt.Errorf("Duplicate template:%s %s", t.Path(), t.Source())
	}

	return err
}

// ParseString parses a string template
func (t *JSONTemplate) ParseString(s string) error {
	mu.Lock()
	defer mu.Unlock()

	err := t.BaseTemplate.ParseString(s)

	// Add to our template set
	if jsonTemplateSet.Lookup(t.Path()) == nil {
		_, err = jsonTemplateSet.New(t.path).Parse(t.Source())
	} else {
		err = fmt.Errorf("Duplicate template:%s %s", t.Path(), t.Source())
	}

	return err
}

// Finalize the template set, called after parsing is complete
func (t *JSONTemplate) Finalize(templates map[string]Template) error {

	// Go html/template records dependencies both ways (child <-> parent)
	// tmpl.Templates() includes tmpl and children and parents
	// we only want includes listed as dependencies
	// so just do a simple search of parsed source instead

	// Search source for {{\s template "|`xxx`|" x }} pattern
	paths := templateInclude.FindAllStringSubmatch(t.Source(), -1)

	// For all includes found, add the template to our dependency list
	for _, p := range paths {
		d := templates[p[1]]
		if d != nil {
			t.dependencies = append(t.dependencies, d)
		}
	}

	return nil
}

// Render the template
func (t *JSONTemplate) Render(writer io.Writer, context map[string]interface{}) error {
	jsonMu.RLock()
	defer jsonMu.RUnlock()
	tmpl := jsonTemplateSet.Lookup(t.Path())
	return tmpl.Execute(writer, context)
}
