package infrastructure

import (
	"encoding/json"
	"html/template"
	"io"
)

// GoTemplateRenderer implements TemplateRenderer using html/template
type GoTemplateRenderer struct {
	tmpl *template.Template
}

// ParseFiles parses the named files into a template
func (tr *GoTemplateRenderer) ParseFiles(filenames ...string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"toJson": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
	}
	var err error
	tr.tmpl, err = template.New("").Funcs(funcMap).ParseFiles(filenames...)
	return tr.tmpl, err
}

// ExecuteTemplate applies the template associated with t that has the given name to the specified data object
func (tr *GoTemplateRenderer) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return tr.tmpl.ExecuteTemplate(wr, name, data)
}
