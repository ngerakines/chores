package chores

import (
	"io"
	"text/template"
)

type AssetFunc func(string) ([]byte, error)

type FuncMap template.FuncMap

// Template is neat.
type Template struct {
	AssetFunc AssetFunc
	tmpl      *template.Template
}

// NewTemplate is neat.
func NewTemplate(name string, fn AssetFunc) *Template {
	return &Template{fn, template.New(name)}
}

// Name is neat.
func (t *Template) Name() string {
	return t.tmpl.Name()
}

// Funcs is neat.
func (t *Template) Funcs(funcMap FuncMap) *Template {
	return t.replaceTmpl(t.tmpl.Funcs(template.FuncMap(funcMap)))
}

// ParseFiles is neat.
func (t *Template) ParseFiles(filenames ...string) (*Template, error) {
	fileBytes := []byte{}
	for _, filename := range filenames {
		tmplBytes, err := t.file(filename)
		if err != nil {
			return nil, err
		}
		fileBytes = append(fileBytes, tmplBytes...)
	}
	newTmpl, err := t.tmpl.Parse(string(fileBytes))
	if err != nil {
		return nil, err
	}
	return t.replaceTmpl(newTmpl), nil
}

// Execute is neat.
func (t *Template) Execute(w io.Writer, data interface{}) error {
	return t.tmpl.Execute(w, data)
}

// ExecuteTemplate is neat.
func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return t.tmpl.ExecuteTemplate(wr, name, data)
}

func (t *Template) replaceTmpl(tmpl *template.Template) *Template {
	t.tmpl = tmpl
	return t
}

func (t *Template) file(fileName string) ([]byte, error) {
	tmplBytes, err := t.AssetFunc(fileName)
	if err != nil {
		return nil, err
	}
	return tmplBytes, nil
}
