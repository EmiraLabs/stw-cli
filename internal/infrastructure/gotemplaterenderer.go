package infrastructure

import (
	"html/template"
	"io"
)

// GoTemplateRenderer implements TemplateRenderer using html/template
type GoTemplateRenderer struct {
	tmpl *template.Template
}

// ParseFiles parses the named files into a template
func (tr *GoTemplateRenderer) ParseFiles(filenames ...string) (*template.Template, error) {
	var err error
	tr.tmpl, err = template.ParseFiles(filenames...)
	return tr.tmpl, err
}

// ExecuteTemplate applies the template associated with t that has the given name to the specified data object
func (tr *GoTemplateRenderer) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return tr.tmpl.ExecuteTemplate(wr, name, data)
}
