package parser

import (
	"io"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/bravepickle/templar/internal/core"
)

type TemplateBuilder struct {
	Name     string
	Vars     core.Params
	Template string
	funcMap  template.FuncMap
}

func (t *TemplateBuilder) Build(w io.Writer) error {
	tpl, err := template.New(t.Name).Funcs(t.funcMap).Parse(t.Template)
	if err != nil {
		return err
	}

	return tpl.Execute(w, t.Vars)
}

func NewTemplate(name string, tpl string, vars core.Params) *TemplateBuilder {
	return &TemplateBuilder{
		Name:     name,
		Vars:     vars,
		Template: tpl,
		funcMap:  sprig.TxtFuncMap(),
	}
}
