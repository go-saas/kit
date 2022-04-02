// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package v1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func IsDuplicateTenantName(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DUPLICATE_TENANT_NAME.String() && e.Code == 400
}

func ErrorDuplicateTenantName(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_DUPLICATE_TENANT_NAME.String(), fmt.Sprintf(format, args...))
}

func IsTenantNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_TENANT_NOT_FOUND.String() && e.Code == 404
}

func ErrorTenantNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(404, ErrorReason_TENANT_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

func IsTenantForbidden(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_TENANT_FORBIDDEN.String() && e.Code == 403
}

func ErrorTenantForbidden(format string, args ...interface{}) *errors.Error {
	return errors.New(403, ErrorReason_TENANT_FORBIDDEN.String(), fmt.Sprintf(format, args...))
}

func IsTenantNotReady(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_TENANT_NOT_READY.String() && e.Code == 403
}

func ErrorTenantNotReady(format string, args ...interface{}) *errors.Error {
	return errors.New(403, ErrorReason_TENANT_NOT_READY.String(), fmt.Sprintf(format, args...))
}
