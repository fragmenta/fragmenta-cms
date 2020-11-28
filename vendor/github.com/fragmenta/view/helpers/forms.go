package helpers

import (
	"fmt"
	got "html/template"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// FORMS

// These should probably use templates from or from lib, so that users can change what form fields get generated
// and use templ rather than fmt.Sprintf

// We need to set this token in the session on the get request for the form

// CSRF generates an input field tag containing a CSRF token
func CSRF() got.HTML {
	token := "my_csrf_token" // instead of generating this here, should we instead get router or app to generate and put into the context?
	output := fmt.Sprintf("<input type='hidden' name='csrf' value='%s'>", token)
	return got.HTML(output)
}

// Field accepts name string, value interface{}, fieldType string, args ...string
func Field(label string, name string, v interface{}, args ...string) got.HTML {
	attributes := ""
	if len(args) > 0 {
		attributes = strings.Join(args, " ")
	}
	// If no type, add it to attributes
	if !strings.Contains(attributes, "type=") {
		attributes = attributes + " type=\"text\""
	}

	tmpl :=
		`<div class="field">
         <label>%s</label>
         <input name="%s" value="%s" %s>
         </div>`

	if label == "" {
		tmpl = `%s<input name="%s" value="%s" %s>`
	}

	output := fmt.Sprintf(tmpl, Escape(label), Escape(name), Escape(fmt.Sprintf("%v", v)), attributes)

	return got.HTML(output)
}

// DateField sets up a date field with a data-date attribute storing the real date
func DateField(label string, name string, t time.Time, args ...string) got.HTML {

	// NB we use text type for date fields because of inconsistent browser behaviour
	// and to support our own date picker popups
	tmpl :=
		`<div class="field">
         <label>%s</label>
         <input name="%s" id="%s" class="date_field" type="text" value="%s" data-date="%s" %s autocomplete="off">
         </div>`

	attributes := ""
	if len(args) > 0 {
		attributes = strings.Join(args, " ")
	}
	output := fmt.Sprintf(tmpl, Escape(label), Escape(name), Escape(name), Date(t), Date(t, "2006-01-02"), attributes)

	return got.HTML(output)
}

// TextArea returns a field div containing a textarea
func TextArea(label string, name string, v interface{}, args ...string) got.HTML {
	attributes := ""
	if len(args) > 0 {
		attributes = strings.Join(args, " ")
	}

	fieldTemplate :=
		`<div class="field">
       <label>%s</label>
       <textarea name="%s" %s>%v</textarea>
       </div>`
	output := fmt.Sprintf(fieldTemplate,
		Escape(label),
		Escape(name),
		attributes, // NB we do not escape attributes, which may contain HTML
		v)          // NB value may contain HTML

	return got.HTML(output)
}

// TODO flip the select helpers to use Selectable all the time?
// Redefine concrete type Option as a Selectable and this should be doable?

// Selectable provides an interface for options in a select
type Selectable interface {
	SelectName() string
	SelectValue() string
}

// SelectableOption provides a concrete implementation of Selectable - this should be called string option or similar
type SelectableOption struct {
	Name  string
	Value string
}

// SelectName returns the public name for this select option
func (o SelectableOption) SelectName() string {
	return o.Name
}

// SelectValue returns the value for this select option
func (o SelectableOption) SelectValue() string {
	return o.Value
}

// StringOptions creates an array of selectables from strings
func StringOptions(args ...string) []SelectableOption {
	var options []SelectableOption
	// Construct a slice of options from these strings

	for _, s := range args {
		options = append(options, SelectableOption{s, s})
	}

	return options
}

// NumberOptions creates an array of selectables, with an optional min and max value supplied as arguments
func NumberOptions(args ...int64) []SelectableOption {

	min := int64(0)
	max := int64(50)

	if len(args) > 0 {
		min = args[0]
	}

	if len(args) > 1 {
		max = args[1]
	}

	var options []SelectableOption

	for i := min; i <= max; i++ {
		v := strconv.Itoa(int(i))
		n := v

		options = append(options, SelectableOption{n, v})
	}

	return options
}

// Better to use an interface and not reflect here - Would rather avoid use of reflect...

// OptionsForSelect creates a select field given an array of keys and values in order
func OptionsForSelect(value interface{}, options interface{}) got.HTML {

	stringValue := fmt.Sprintf("%v", value)

	output := ""

	switch reflect.TypeOf(options).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(options)
		for i := 0; i < s.Len(); i++ {
			o := s.Index(i).Interface().(Selectable)
			sel := ""
			if o.SelectValue() == stringValue {
				sel = "selected"
			}

			output += fmt.Sprintf(`<option value="%s" %s>%s</option>
`, o.SelectValue(), sel, Escape(o.SelectName()))

		}
	}

	return got.HTML(output)
}

// SelectArray creates a select field given an array of keys and values in order
func SelectArray(label string, name string, value interface{}, options interface{}) got.HTML {

	stringValue := fmt.Sprintf("%v", value)

	tmpl :=
		`<div class="field">
      <label>%s</label>
      <select type="select" name="%s" id="%s">
      %s
      </select>
      </div>`

	if label == "" {
		tmpl = `%s<select type="select" name="%s" id="%s">
%s
</select>`
	}

	opts := ""

	switch reflect.TypeOf(options).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(options)
		for i := 0; i < s.Len(); i++ {
			o := s.Index(i).Interface().(Selectable)
			sel := ""
			if o.SelectValue() == stringValue {
				sel = "selected"
			}

			opts += fmt.Sprintf(`<option value="%s" %s>%s</option>
`, o.SelectValue(), sel, Escape(o.SelectName()))

		}
	}

	output := fmt.Sprintf(tmpl, Escape(label), Escape(name), Escape(name), opts)

	return got.HTML(output)
}

// FIXME - make Option conform to Selectable interface and use that instead of concrete type below

// Option type contains number and string
type Option struct {
	Id   int64  // The value - FIXME migrate to ID and use as interface
	Name string // The name
}

// SelectName returns the public name for this select option
func (o Option) SelectName() string {
	return o.Name
}

// SelectValue returns the value for this select option
func (o Option) SelectValue() string {
	return fmt.Sprintf("%d", o.Id)
}

// SelectID returns the value for this select option as an int64
func (o Option) SelectID() int64 {
	return o.Id
}

// ID returns the value for this select option as an int64
// this is supplied in case of use in templates
func (o Option) ID() int64 {
	return o.Id
}

// SetID is deprecated
func (o Option) SetID(id int64) {
	o.Id = id
}

// NewOption returns a new option, given an id and name
func NewOption(id int64, name string) Option {
	return Option{Id: id, Name: name}
}

// Select creates a select field given an array of keys and values in order
func Select(label string, name string, value int64, options []Option) got.HTML {

	tmpl :=
		`<div class="field">
      <label>%s</label>
      <select type="select" name="%s">
      %s
      </select>
      </div>`

	if label == "" {
		tmpl = `%s<select type="select" name="%s">
      %s
      </select>`
	}

	opts := ""
	for _, o := range options {

		s := ""
		if o.Id == value {
			s = "selected"
		}

		opts += fmt.Sprintf(`<option value="%d" %s>%s</option>
`, o.Id, s, Escape(o.Name))
	}

	output := fmt.Sprintf(tmpl, Escape(label), Escape(name), opts)

	return got.HTML(output)
}
