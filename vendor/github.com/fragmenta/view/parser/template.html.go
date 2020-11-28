package parser

import (
	"fmt"
	got "html/template"
	"io"
	"sync"
)

var (
	mu              sync.RWMutex  // Shared mutex to go with shared template set, because of dev reloads
	htmlTemplateSet *got.Template // This is a shared template set for HTML templates
)

// HTMLTemplate represents an HTML template using go HTML/template
type HTMLTemplate struct {
	BaseTemplate
}

// Setup performs setup before parsing templates
func (t *HTMLTemplate) Setup(helpers FuncMap) error {
	mu.Lock()
	defer mu.Unlock()
	htmlTemplateSet = got.New("").Funcs(got.FuncMap(helpers))
	return nil
}

// CanParseFile returns true if this parser handles this path
func (t *HTMLTemplate) CanParseFile(path string) bool {
	allowed := []string{".html.got", ".xml.got"}
	return suffixes(path, allowed)
}

// NewTemplate returns a new template for this type
func (t *HTMLTemplate) NewTemplate(fullpath, path string) (Template, error) {
	template := new(HTMLTemplate)
	template.fullpath = fullpath
	template.path = path
	return template, nil
}

// Parse the template at path
func (t *HTMLTemplate) Parse() error {
	mu.Lock()
	defer mu.Unlock()
	err := t.BaseTemplate.Parse()
	if err != nil {
		return err
	}

	// Add to our template set - NB duplicates not allowed by golang templates
	if htmlTemplateSet.Lookup(t.Path()) == nil {
		_, err = htmlTemplateSet.New(t.path).Parse(t.Source())
	} else {
		err = fmt.Errorf("Duplicate template:%s %s", t.Path(), t.Source())
	}

	return err
}

// ParseString parses a string template
func (t *HTMLTemplate) ParseString(s string) error {
	mu.Lock()
	defer mu.Unlock()
	err := t.BaseTemplate.ParseString(s)

	// Add to our template set
	if htmlTemplateSet.Lookup(t.Path()) == nil {
		_, err = htmlTemplateSet.New(t.path).Parse(t.Source())
	} else {
		err = fmt.Errorf("Duplicate template:%s %s", t.Path(), t.Source())
	}

	return err
}

// Finalize the template set, called after parsing is complete
func (t *HTMLTemplate) Finalize(templates map[string]Template) error {

	// Go html/template records dependencies both ways (child <-> parent)
	// tmpl.Templates() includes tmpl and children and parents
	// we only want includes listed as dependencies
	// so just do a simple search of unparsed source instead

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

// Render the template to the given writer, returning an error
func (t *HTMLTemplate) Render(writer io.Writer, context map[string]interface{}) error {
	mu.RLock()
	defer mu.RUnlock()
	tmpl := htmlTemplateSet.Lookup(t.Path())
	if tmpl == nil {
		return fmt.Errorf("#error loading template for %s", t.Path())
	}
	return tmpl.Execute(writer, context)
}
