{{ range .Errors }}

{{ if .HasComment }}{{ .Comment }}{{ end -}}
func Is{{.CamelValue}}(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == {{ .Name }}_{{ .Value }}.String() && e.Code == {{ .HTTPCode }}
}

{{ if .HasComment }}{{ .Comment }}{{ end -}}
func Error{{.CamelValue}}Localized(ctx context.Context, data map[string]interface{}, pluralCount interface{}) *errors.Error {
    localizer := localize.FromContext(ctx)
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
