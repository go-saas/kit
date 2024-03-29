// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: product/private/biz/job.proto

package biz

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on ProductUpdatedJobParam with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ProductUpdatedJobParam) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ProductUpdatedJobParam with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ProductUpdatedJobParamMultiError, or nil if none found.
func (m *ProductUpdatedJobParam) ValidateAll() error {
	return m.validate(true)
}

func (m *ProductUpdatedJobParam) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ProductId

	// no validation rules for ProductVersion

	// no validation rules for TenantId

	// no validation rules for IsDelete

	for idx, item := range m.GetSyncLinks() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ProductUpdatedJobParamValidationError{
						field:  fmt.Sprintf("SyncLinks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ProductUpdatedJobParamValidationError{
						field:  fmt.Sprintf("SyncLinks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ProductUpdatedJobParamValidationError{
					field:  fmt.Sprintf("SyncLinks[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ProductUpdatedJobParamMultiError(errors)
	}

	return nil
}

// ProductUpdatedJobParamMultiError is an error wrapping multiple validation
// errors returned by ProductUpdatedJobParam.ValidateAll() if the designated
// constraints aren't met.
type ProductUpdatedJobParamMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ProductUpdatedJobParamMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ProductUpdatedJobParamMultiError) AllErrors() []error { return m }

// ProductUpdatedJobParamValidationError is the validation error returned by
// ProductUpdatedJobParam.Validate if the designated constraints aren't met.
type ProductUpdatedJobParamValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ProductUpdatedJobParamValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ProductUpdatedJobParamValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ProductUpdatedJobParamValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ProductUpdatedJobParamValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ProductUpdatedJobParamValidationError) ErrorName() string {
	return "ProductUpdatedJobParamValidationError"
}

// Error satisfies the builtin error interface
func (e ProductUpdatedJobParamValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sProductUpdatedJobParam.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ProductUpdatedJobParamValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ProductUpdatedJobParamValidationError{}

// Validate checks the field values on ProductUpdatedJobParam_SyncLink with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ProductUpdatedJobParam_SyncLink) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ProductUpdatedJobParam_SyncLink with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// ProductUpdatedJobParam_SyncLinkMultiError, or nil if none found.
func (m *ProductUpdatedJobParam_SyncLink) ValidateAll() error {
	return m.validate(true)
}

func (m *ProductUpdatedJobParam_SyncLink) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ProviderName

	// no validation rules for ProviderId

	if len(errors) > 0 {
		return ProductUpdatedJobParam_SyncLinkMultiError(errors)
	}

	return nil
}

// ProductUpdatedJobParam_SyncLinkMultiError is an error wrapping multiple
// validation errors returned by ProductUpdatedJobParam_SyncLink.ValidateAll()
// if the designated constraints aren't met.
type ProductUpdatedJobParam_SyncLinkMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ProductUpdatedJobParam_SyncLinkMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ProductUpdatedJobParam_SyncLinkMultiError) AllErrors() []error { return m }

// ProductUpdatedJobParam_SyncLinkValidationError is the validation error
// returned by ProductUpdatedJobParam_SyncLink.Validate if the designated
// constraints aren't met.
type ProductUpdatedJobParam_SyncLinkValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ProductUpdatedJobParam_SyncLinkValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ProductUpdatedJobParam_SyncLinkValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ProductUpdatedJobParam_SyncLinkValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ProductUpdatedJobParam_SyncLinkValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ProductUpdatedJobParam_SyncLinkValidationError) ErrorName() string {
	return "ProductUpdatedJobParam_SyncLinkValidationError"
}

// Error satisfies the builtin error interface
func (e ProductUpdatedJobParam_SyncLinkValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sProductUpdatedJobParam_SyncLink.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ProductUpdatedJobParam_SyncLinkValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ProductUpdatedJobParam_SyncLinkValidationError{}
