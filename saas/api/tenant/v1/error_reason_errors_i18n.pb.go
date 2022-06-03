// Code generated by protoc-gen-go-errors-i18n. DO NOT EDIT.

package v1

import (
	errors "github.com/go-kratos/kratos/v2/errors"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func ErrorDuplicateTenantNameLocalized(localizer *i18n.Localizer, data map[string]interface{}, pluralCount interface{}) *errors.Error {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "DuplicateTenantName",
		},
		TemplateData: data,
		PluralCount:  pluralCount,
	})
	if err == nil {
		return errors.New(400, ErrorReason_DUPLICATE_TENANT_NAME.String(), msg)
	} else {
		return errors.New(400, ErrorReason_DUPLICATE_TENANT_NAME.String(), "")
	}

}

func ErrorTenantNotFoundLocalized(localizer *i18n.Localizer, data map[string]interface{}, pluralCount interface{}) *errors.Error {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "TenantNotFound",
		},
		TemplateData: data,
		PluralCount:  pluralCount,
	})
	if err == nil {
		return errors.New(404, ErrorReason_TENANT_NOT_FOUND.String(), msg)
	} else {
		return errors.New(404, ErrorReason_TENANT_NOT_FOUND.String(), "")
	}

}

func ErrorTenantForbiddenLocalized(localizer *i18n.Localizer, data map[string]interface{}, pluralCount interface{}) *errors.Error {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "TenantForbidden",
		},
		TemplateData: data,
		PluralCount:  pluralCount,
	})
	if err == nil {
		return errors.New(403, ErrorReason_TENANT_FORBIDDEN.String(), msg)
	} else {
		return errors.New(403, ErrorReason_TENANT_FORBIDDEN.String(), "")
	}

}

func ErrorTenantNotReadyLocalized(localizer *i18n.Localizer, data map[string]interface{}, pluralCount interface{}) *errors.Error {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "TenantNotReady",
		},
		TemplateData: data,
		PluralCount:  pluralCount,
	})
	if err == nil {
		return errors.New(403, ErrorReason_TENANT_NOT_READY.String(), msg)
	} else {
		return errors.New(403, ErrorReason_TENANT_NOT_READY.String(), "")
	}

}