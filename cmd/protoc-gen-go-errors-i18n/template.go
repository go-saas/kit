package main

import (
	"bytes"
	"text/template"
)

var errorsTemplate = `
{{ range .Errors }}

func Error{{.CamelValue}}Localized(localizer *i18n.Localizer, data map[string]interface{}, pluralCount interface{}) *errors.Error {
     if localizer == nil {
		return errors.New({{.HTTPCode}}, {{.Name}}_{{.Value}}.String(), "")
	 }
     msg, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "{{.MsgKey}}",
		},
		TemplateData: data,
		PluralCount: pluralCount,
	})
	if err == nil {
		return errors.New({{.HTTPCode}}, {{.Name}}_{{.Value}}.String(), msg)
	} else {
		return errors.New({{.HTTPCode}}, {{.Name}}_{{.Value}}.String(), "")
	}
	
}

{{- end }}
`

type errorInfo struct {
	Name       string
	Value      string
	HTTPCode   int
	CamelValue string
	MsgKey     string
}

type errorWrapper struct {
	Errors []*errorInfo
}

func (e *errorWrapper) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errors").Parse(errorsTemplate)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return buf.String()
}
