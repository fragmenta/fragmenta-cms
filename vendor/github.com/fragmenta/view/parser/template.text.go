package parser

import (
	"fmt"
	"io"
	got "text/template"
)

var textTemplateSet *got.Template

// TextTemplate using go text/template
type TextTemplate struct {
	BaseTemplate
}

// Setup runs before parsing templates
func (t *TextTemplate) Setup(helpers FuncMap) error {
	textTemplateSet = got.New("").Funcs(got.FuncMap(helpers))
	return nil
}

// CanParseFile returns true if this parser handles this file path?
func (t *TextTemplate) CanParseFile(path string) bool {
	allowed := []string{".text.got", ".csv.got"}
	return suffixes(path, allowed)
}

// NewTemplate returns a new template of this type
func (t *TextTemplate) NewTemplate(fullpath, path string) (Template, error) {
	template := new(TextTemplate)
	template.fullpath = fullpath
	template.path = path
	return template, nil
}

// Parse the template
func (t *TextTemplate) Parse() error {
	err := t.BaseTemplate.Parse()

	// Add to our template set
	if textTemplateSet.Lookup(t.path) == nil {
		_, err = textTemplateSet.New(t.path).Parse(t.Source())
	} else {
		err = fmt.Errorf("Duplicate template:%s %s", t.Path(), t.Source())
	}

	return err
}

// ParseString a string template
func (t *TextTemplate) ParseString(s string) error {
	err := t.BaseTemplate.ParseString(s)

	// Add to our template set
	if textTemplateSet.Lookup(t.Path()) == nil {
		_, err = textTemplateSet.New(t.path).Parse(t.Source())
	} else {
		err = fmt.Errorf("Duplicate template:%s %s", t.Path(), t.Source())
	}

	return err
}

// Finalize the template set, called after parsing is complete
// Record a list of dependent templates (for breaking caches automatically)
func (t *TextTemplate) Finalize(templates map[string]Template) error {

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

// Render renders the template
func (t *TextTemplate) Render(writer io.Writer, context map[string]interface{}) error {
	tmpl := t.goTemplate()
	if tmpl == nil {
		return fmt.Errorf("Error rendering template:%s %s", t.Path(), t.Source())
	}

	return tmpl.Execute(writer, context)
}

// goTemplate returns teh underlying go template
func (t *TextTemplate) goTemplate() *got.Template {
	return textTemplateSet.Lookup(t.Path())
}
