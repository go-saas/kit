// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: user/api/permission/v1/permission.proto

package v1

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

// Validate checks the field values on GetCurrentPermissionRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetCurrentPermissionRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetCurrentPermissionRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetCurrentPermissionRequestMultiError, or nil if none found.
func (m *GetCurrentPermissionRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetCurrentPermissionRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return GetCurrentPermissionRequestMultiError(errors)
	}
	return nil
}

// GetCurrentPermissionRequestMultiError is an error wrapping multiple
// validation errors returned by GetCurrentPermissionRequest.ValidateAll() if
// the designated constraints aren't met.
type GetCurrentPermissionRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetCurrentPermissionRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetCurrentPermissionRequestMultiError) AllErrors() []error { return m }

// GetCurrentPermissionRequestValidationError is the validation error returned
// by GetCurrentPermissionRequest.Validate if the designated constraints
// aren't met.
type GetCurrentPermissionRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCurrentPermissionRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCurrentPermissionRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCurrentPermissionRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCurrentPermissionRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCurrentPermissionRequestValidationError) ErrorName() string {
	return "GetCurrentPermissionRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetCurrentPermissionRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCurrentPermissionRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCurrentPermissionRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCurrentPermissionRequestValidationError{}

// Validate checks the field values on GetCurrentPermissionReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetCurrentPermissionReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetCurrentPermissionReply with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetCurrentPermissionReplyMultiError, or nil if none found.
func (m *GetCurrentPermissionReply) ValidateAll() error {
	return m.validate(true)
}

func (m *GetCurrentPermissionReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetAcl() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetCurrentPermissionReplyValidationError{
						field:  fmt.Sprintf("Acl[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetCurrentPermissionReplyValidationError{
						field:  fmt.Sprintf("Acl[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetCurrentPermissionReplyValidationError{
					field:  fmt.Sprintf("Acl[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetCurrentPermissionReplyMultiError(errors)
	}
	return nil
}

// GetCurrentPermissionReplyMultiError is an error wrapping multiple validation
// errors returned by GetCurrentPermissionReply.ValidateAll() if the
// designated constraints aren't met.
type GetCurrentPermissionReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetCurrentPermissionReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetCurrentPermissionReplyMultiError) AllErrors() []error { return m }

// GetCurrentPermissionReplyValidationError is the validation error returned by
// GetCurrentPermissionReply.Validate if the designated constraints aren't met.
type GetCurrentPermissionReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCurrentPermissionReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCurrentPermissionReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCurrentPermissionReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCurrentPermissionReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCurrentPermissionReplyValidationError) ErrorName() string {
	return "GetCurrentPermissionReplyValidationError"
}

// Error satisfies the builtin error interface
func (e GetCurrentPermissionReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCurrentPermissionReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCurrentPermissionReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCurrentPermissionReplyValidationError{}

// Validate checks the field values on CheckPermissionRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CheckPermissionRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CheckPermissionRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CheckPermissionRequestMultiError, or nil if none found.
func (m *CheckPermissionRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CheckPermissionRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Namespace

	// no validation rules for Resource

	// no validation rules for Action

	if len(errors) > 0 {
		return CheckPermissionRequestMultiError(errors)
	}
	return nil
}

// CheckPermissionRequestMultiError is an error wrapping multiple validation
// errors returned by CheckPermissionRequest.ValidateAll() if the designated
// constraints aren't met.
type CheckPermissionRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CheckPermissionRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CheckPermissionRequestMultiError) AllErrors() []error { return m }

// CheckPermissionRequestValidationError is the validation error returned by
// CheckPermissionRequest.Validate if the designated constraints aren't met.
type CheckPermissionRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CheckPermissionRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CheckPermissionRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CheckPermissionRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CheckPermissionRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CheckPermissionRequestValidationError) ErrorName() string {
	return "CheckPermissionRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CheckPermissionRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCheckPermissionRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CheckPermissionRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CheckPermissionRequestValidationError{}

// Validate checks the field values on CheckPermissionReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CheckPermissionReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CheckPermissionReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CheckPermissionReplyMultiError, or nil if none found.
func (m *CheckPermissionReply) ValidateAll() error {
	return m.validate(true)
}

func (m *CheckPermissionReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Effect

	if len(errors) > 0 {
		return CheckPermissionReplyMultiError(errors)
	}
	return nil
}

// CheckPermissionReplyMultiError is an error wrapping multiple validation
// errors returned by CheckPermissionReply.ValidateAll() if the designated
// constraints aren't met.
type CheckPermissionReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CheckPermissionReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CheckPermissionReplyMultiError) AllErrors() []error { return m }

// CheckPermissionReplyValidationError is the validation error returned by
// CheckPermissionReply.Validate if the designated constraints aren't met.
type CheckPermissionReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CheckPermissionReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CheckPermissionReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CheckPermissionReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CheckPermissionReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CheckPermissionReplyValidationError) ErrorName() string {
	return "CheckPermissionReplyValidationError"
}

// Error satisfies the builtin error interface
func (e CheckPermissionReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCheckPermissionReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CheckPermissionReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CheckPermissionReplyValidationError{}

// Validate checks the field values on CheckSubjectsPermissionRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CheckSubjectsPermissionRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CheckSubjectsPermissionRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// CheckSubjectsPermissionRequestMultiError, or nil if none found.
func (m *CheckSubjectsPermissionRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CheckSubjectsPermissionRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Namespace

	// no validation rules for Resource

	// no validation rules for Action

	// no validation rules for TenantId

	if len(errors) > 0 {
		return CheckSubjectsPermissionRequestMultiError(errors)
	}
	return nil
}

// CheckSubjectsPermissionRequestMultiError is an error wrapping multiple
// validation errors returned by CheckSubjectsPermissionRequest.ValidateAll()
// if the designated constraints aren't met.
type CheckSubjectsPermissionRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CheckSubjectsPermissionRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CheckSubjectsPermissionRequestMultiError) AllErrors() []error { return m }

// CheckSubjectsPermissionRequestValidationError is the validation error
// returned by CheckSubjectsPermissionRequest.Validate if the designated
// constraints aren't met.
type CheckSubjectsPermissionRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CheckSubjectsPermissionRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CheckSubjectsPermissionRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CheckSubjectsPermissionRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CheckSubjectsPermissionRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CheckSubjectsPermissionRequestValidationError) ErrorName() string {
	return "CheckSubjectsPermissionRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CheckSubjectsPermissionRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCheckSubjectsPermissionRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CheckSubjectsPermissionRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CheckSubjectsPermissionRequestValidationError{}

// Validate checks the field values on CheckSubjectsPermissionReply with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CheckSubjectsPermissionReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CheckSubjectsPermissionReply with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CheckSubjectsPermissionReplyMultiError, or nil if none found.
func (m *CheckSubjectsPermissionReply) ValidateAll() error {
	return m.validate(true)
}

func (m *CheckSubjectsPermissionReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Effect

	if len(errors) > 0 {
		return CheckSubjectsPermissionReplyMultiError(errors)
	}
	return nil
}

// CheckSubjectsPermissionReplyMultiError is an error wrapping multiple
// validation errors returned by CheckSubjectsPermissionReply.ValidateAll() if
// the designated constraints aren't met.
type CheckSubjectsPermissionReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CheckSubjectsPermissionReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CheckSubjectsPermissionReplyMultiError) AllErrors() []error { return m }

// CheckSubjectsPermissionReplyValidationError is the validation error returned
// by CheckSubjectsPermissionReply.Validate if the designated constraints
// aren't met.
type CheckSubjectsPermissionReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CheckSubjectsPermissionReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CheckSubjectsPermissionReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CheckSubjectsPermissionReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CheckSubjectsPermissionReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CheckSubjectsPermissionReplyValidationError) ErrorName() string {
	return "CheckSubjectsPermissionReplyValidationError"
}

// Error satisfies the builtin error interface
func (e CheckSubjectsPermissionReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCheckSubjectsPermissionReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CheckSubjectsPermissionReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CheckSubjectsPermissionReplyValidationError{}

// Validate checks the field values on Permission with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Permission) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Permission with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PermissionMultiError, or
// nil if none found.
func (m *Permission) ValidateAll() error {
	return m.validate(true)
}

func (m *Permission) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Namespace

	// no validation rules for Resource

	// no validation rules for Action

	// no validation rules for Subject

	// no validation rules for Effect

	if len(errors) > 0 {
		return PermissionMultiError(errors)
	}
	return nil
}

// PermissionMultiError is an error wrapping multiple validation errors
// returned by Permission.ValidateAll() if the designated constraints aren't met.
type PermissionMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PermissionMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PermissionMultiError) AllErrors() []error { return m }

// PermissionValidationError is the validation error returned by
// Permission.Validate if the designated constraints aren't met.
type PermissionValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PermissionValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PermissionValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PermissionValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PermissionValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PermissionValidationError) ErrorName() string { return "PermissionValidationError" }

// Error satisfies the builtin error interface
func (e PermissionValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPermission.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PermissionValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PermissionValidationError{}

// Validate checks the field values on UpdateSubjectPermissionRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateSubjectPermissionRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateSubjectPermissionRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// UpdateSubjectPermissionRequestMultiError, or nil if none found.
func (m *UpdateSubjectPermissionRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateSubjectPermissionRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Subject

	for idx, item := range m.GetAcl() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UpdateSubjectPermissionRequestValidationError{
						field:  fmt.Sprintf("Acl[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UpdateSubjectPermissionRequestValidationError{
						field:  fmt.Sprintf("Acl[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UpdateSubjectPermissionRequestValidationError{
					field:  fmt.Sprintf("Acl[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return UpdateSubjectPermissionRequestMultiError(errors)
	}
	return nil
}

// UpdateSubjectPermissionRequestMultiError is an error wrapping multiple
// validation errors returned by UpdateSubjectPermissionRequest.ValidateAll()
// if the designated constraints aren't met.
type UpdateSubjectPermissionRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateSubjectPermissionRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateSubjectPermissionRequestMultiError) AllErrors() []error { return m }

// UpdateSubjectPermissionRequestValidationError is the validation error
// returned by UpdateSubjectPermissionRequest.Validate if the designated
// constraints aren't met.
type UpdateSubjectPermissionRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateSubjectPermissionRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateSubjectPermissionRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateSubjectPermissionRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateSubjectPermissionRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateSubjectPermissionRequestValidationError) ErrorName() string {
	return "UpdateSubjectPermissionRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateSubjectPermissionRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateSubjectPermissionRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateSubjectPermissionRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateSubjectPermissionRequestValidationError{}

// Validate checks the field values on UpdateSubjectPermissionAcl with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateSubjectPermissionAcl) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateSubjectPermissionAcl with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateSubjectPermissionAclMultiError, or nil if none found.
func (m *UpdateSubjectPermissionAcl) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateSubjectPermissionAcl) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Namespace

	// no validation rules for Resource

	// no validation rules for Action

	// no validation rules for Effect

	if len(errors) > 0 {
		return UpdateSubjectPermissionAclMultiError(errors)
	}
	return nil
}

// UpdateSubjectPermissionAclMultiError is an error wrapping multiple
// validation errors returned by UpdateSubjectPermissionAcl.ValidateAll() if
// the designated constraints aren't met.
type UpdateSubjectPermissionAclMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateSubjectPermissionAclMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateSubjectPermissionAclMultiError) AllErrors() []error { return m }

// UpdateSubjectPermissionAclValidationError is the validation error returned
// by UpdateSubjectPermissionAcl.Validate if the designated constraints aren't met.
type UpdateSubjectPermissionAclValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateSubjectPermissionAclValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateSubjectPermissionAclValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateSubjectPermissionAclValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateSubjectPermissionAclValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateSubjectPermissionAclValidationError) ErrorName() string {
	return "UpdateSubjectPermissionAclValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateSubjectPermissionAclValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateSubjectPermissionAcl.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateSubjectPermissionAclValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateSubjectPermissionAclValidationError{}

// Validate checks the field values on UpdateSubjectPermissionResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateSubjectPermissionResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateSubjectPermissionResponse with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// UpdateSubjectPermissionResponseMultiError, or nil if none found.
func (m *UpdateSubjectPermissionResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateSubjectPermissionResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetAcl() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UpdateSubjectPermissionResponseValidationError{
						field:  fmt.Sprintf("Acl[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UpdateSubjectPermissionResponseValidationError{
						field:  fmt.Sprintf("Acl[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UpdateSubjectPermissionResponseValidationError{
					field:  fmt.Sprintf("Acl[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return UpdateSubjectPermissionResponseMultiError(errors)
	}
	return nil
}

// UpdateSubjectPermissionResponseMultiError is an error wrapping multiple
// validation errors returned by UpdateSubjectPermissionResponse.ValidateAll()
// if the designated constraints aren't met.
type UpdateSubjectPermissionResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateSubjectPermissionResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateSubjectPermissionResponseMultiError) AllErrors() []error { return m }

// UpdateSubjectPermissionResponseValidationError is the validation error
// returned by UpdateSubjectPermissionResponse.Validate if the designated
// constraints aren't met.
type UpdateSubjectPermissionResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateSubjectPermissionResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateSubjectPermissionResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateSubjectPermissionResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateSubjectPermissionResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateSubjectPermissionResponseValidationError) ErrorName() string {
	return "UpdateSubjectPermissionResponseValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateSubjectPermissionResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateSubjectPermissionResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateSubjectPermissionResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateSubjectPermissionResponseValidationError{}
