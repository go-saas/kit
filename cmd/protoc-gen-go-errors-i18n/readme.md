# protoc-gen-go-errors-i18n

generate [i18n](https://github.com/go-saas/go-i18n) method for [errors](https://github.com/go-kratos/kratos/tree/main/cmd/protoc-gen-go-errors)

for example
```go
func ErrorInvalidCredentialsLocalized(localizer *i18n.Localizer, data map[string]interface{}, pluralCount interface{}) *errors.Error {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "InvalidCredentials",
		},
		TemplateData: data,
		PluralCount:  pluralCount,
	})
	if err == nil {
		return errors.New(400, ErrorReason_INVALID_CREDENTIALS.String(), msg)
	} else {
		return errors.New(400, ErrorReason_INVALID_CREDENTIALS.String(), "")
	}

}
```