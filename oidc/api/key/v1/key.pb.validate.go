// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: oidc/api/key/v1/key.proto

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

// Validate checks the field values on DeleteJsonWebKeySetRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteJsonWebKeySetRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteJsonWebKeySetRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteJsonWebKeySetRequestMultiError, or nil if none found.
func (m *DeleteJsonWebKeySetRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteJsonWebKeySetRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetSet()) < 1 {
		err := DeleteJsonWebKeySetRequestValidationError{
			field:  "Set",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DeleteJsonWebKeySetRequestMultiError(errors)
	}

	return nil
}

// DeleteJsonWebKeySetRequestMultiError is an error wrapping multiple
// validation errors returned by DeleteJsonWebKeySetRequest.ValidateAll() if
// the designated constraints aren't met.
type DeleteJsonWebKeySetRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteJsonWebKeySetRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteJsonWebKeySetRequestMultiError) AllErrors() []error { return m }

// DeleteJsonWebKeySetRequestValidationError is the validation error returned
// by DeleteJsonWebKeySetRequest.Validate if the designated constraints aren't met.
type DeleteJsonWebKeySetRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteJsonWebKeySetRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteJsonWebKeySetRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteJsonWebKeySetRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteJsonWebKeySetRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteJsonWebKeySetRequestValidationError) ErrorName() string {
	return "DeleteJsonWebKeySetRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteJsonWebKeySetRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteJsonWebKeySetRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteJsonWebKeySetRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteJsonWebKeySetRequestValidationError{}

// Validate checks the field values on GetJsonWebKeySetRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetJsonWebKeySetRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetJsonWebKeySetRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetJsonWebKeySetRequestMultiError, or nil if none found.
func (m *GetJsonWebKeySetRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetJsonWebKeySetRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetSet()) < 1 {
		err := GetJsonWebKeySetRequestValidationError{
			field:  "Set",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetJsonWebKeySetRequestMultiError(errors)
	}

	return nil
}

// GetJsonWebKeySetRequestMultiError is an error wrapping multiple validation
// errors returned by GetJsonWebKeySetRequest.ValidateAll() if the designated
// constraints aren't met.
type GetJsonWebKeySetRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetJsonWebKeySetRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetJsonWebKeySetRequestMultiError) AllErrors() []error { return m }

// GetJsonWebKeySetRequestValidationError is the validation error returned by
// GetJsonWebKeySetRequest.Validate if the designated constraints aren't met.
type GetJsonWebKeySetRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetJsonWebKeySetRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetJsonWebKeySetRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetJsonWebKeySetRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetJsonWebKeySetRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetJsonWebKeySetRequestValidationError) ErrorName() string {
	return "GetJsonWebKeySetRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetJsonWebKeySetRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetJsonWebKeySetRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetJsonWebKeySetRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetJsonWebKeySetRequestValidationError{}

// Validate checks the field values on UpdateJsonWebKeySetRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateJsonWebKeySetRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateJsonWebKeySetRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateJsonWebKeySetRequestMultiError, or nil if none found.
func (m *UpdateJsonWebKeySetRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateJsonWebKeySetRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetSet()) < 1 {
		err := UpdateJsonWebKeySetRequestValidationError{
			field:  "Set",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetKeys()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UpdateJsonWebKeySetRequestValidationError{
					field:  "Keys",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UpdateJsonWebKeySetRequestValidationError{
					field:  "Keys",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetKeys()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdateJsonWebKeySetRequestValidationError{
				field:  "Keys",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return UpdateJsonWebKeySetRequestMultiError(errors)
	}

	return nil
}

// UpdateJsonWebKeySetRequestMultiError is an error wrapping multiple
// validation errors returned by UpdateJsonWebKeySetRequest.ValidateAll() if
// the designated constraints aren't met.
type UpdateJsonWebKeySetRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateJsonWebKeySetRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateJsonWebKeySetRequestMultiError) AllErrors() []error { return m }

// UpdateJsonWebKeySetRequestValidationError is the validation error returned
// by UpdateJsonWebKeySetRequest.Validate if the designated constraints aren't met.
type UpdateJsonWebKeySetRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateJsonWebKeySetRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateJsonWebKeySetRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateJsonWebKeySetRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateJsonWebKeySetRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateJsonWebKeySetRequestValidationError) ErrorName() string {
	return "UpdateJsonWebKeySetRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateJsonWebKeySetRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateJsonWebKeySetRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateJsonWebKeySetRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateJsonWebKeySetRequestValidationError{}

// Validate checks the field values on CreateJsonWebKeySetRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateJsonWebKeySetRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateJsonWebKeySetRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateJsonWebKeySetRequestMultiError, or nil if none found.
func (m *CreateJsonWebKeySetRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateJsonWebKeySetRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetSet()) < 1 {
		err := CreateJsonWebKeySetRequestValidationError{
			field:  "Set",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetKeys()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreateJsonWebKeySetRequestValidationError{
					field:  "Keys",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreateJsonWebKeySetRequestValidationError{
					field:  "Keys",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetKeys()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateJsonWebKeySetRequestValidationError{
				field:  "Keys",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CreateJsonWebKeySetRequestMultiError(errors)
	}

	return nil
}

// CreateJsonWebKeySetRequestMultiError is an error wrapping multiple
// validation errors returned by CreateJsonWebKeySetRequest.ValidateAll() if
// the designated constraints aren't met.
type CreateJsonWebKeySetRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateJsonWebKeySetRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateJsonWebKeySetRequestMultiError) AllErrors() []error { return m }

// CreateJsonWebKeySetRequestValidationError is the validation error returned
// by CreateJsonWebKeySetRequest.Validate if the designated constraints aren't met.
type CreateJsonWebKeySetRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateJsonWebKeySetRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateJsonWebKeySetRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateJsonWebKeySetRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateJsonWebKeySetRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateJsonWebKeySetRequestValidationError) ErrorName() string {
	return "CreateJsonWebKeySetRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateJsonWebKeySetRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateJsonWebKeySetRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateJsonWebKeySetRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateJsonWebKeySetRequestValidationError{}

// Validate checks the field values on JsonWebKeySetGeneratorRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *JsonWebKeySetGeneratorRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on JsonWebKeySetGeneratorRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// JsonWebKeySetGeneratorRequestMultiError, or nil if none found.
func (m *JsonWebKeySetGeneratorRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *JsonWebKeySetGeneratorRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetAlg()) < 1 {
		err := JsonWebKeySetGeneratorRequestValidationError{
			field:  "Alg",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetKid()) < 1 {
		err := JsonWebKeySetGeneratorRequestValidationError{
			field:  "Kid",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetUse()) < 1 {
		err := JsonWebKeySetGeneratorRequestValidationError{
			field:  "Use",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return JsonWebKeySetGeneratorRequestMultiError(errors)
	}

	return nil
}

// JsonWebKeySetGeneratorRequestMultiError is an error wrapping multiple
// validation errors returned by JsonWebKeySetGeneratorRequest.ValidateAll()
// if the designated constraints aren't met.
type JsonWebKeySetGeneratorRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m JsonWebKeySetGeneratorRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m JsonWebKeySetGeneratorRequestMultiError) AllErrors() []error { return m }

// JsonWebKeySetGeneratorRequestValidationError is the validation error
// returned by JsonWebKeySetGeneratorRequest.Validate if the designated
// constraints aren't met.
type JsonWebKeySetGeneratorRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e JsonWebKeySetGeneratorRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e JsonWebKeySetGeneratorRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e JsonWebKeySetGeneratorRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e JsonWebKeySetGeneratorRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e JsonWebKeySetGeneratorRequestValidationError) ErrorName() string {
	return "JsonWebKeySetGeneratorRequestValidationError"
}

// Error satisfies the builtin error interface
func (e JsonWebKeySetGeneratorRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sJsonWebKeySetGeneratorRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = JsonWebKeySetGeneratorRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = JsonWebKeySetGeneratorRequestValidationError{}

// Validate checks the field values on DeleteJsonWebKeyRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteJsonWebKeyRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteJsonWebKeyRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteJsonWebKeyRequestMultiError, or nil if none found.
func (m *DeleteJsonWebKeyRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteJsonWebKeyRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetSet()) < 1 {
		err := DeleteJsonWebKeyRequestValidationError{
			field:  "Set",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetKid()) < 1 {
		err := DeleteJsonWebKeyRequestValidationError{
			field:  "Kid",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DeleteJsonWebKeyRequestMultiError(errors)
	}

	return nil
}

// DeleteJsonWebKeyRequestMultiError is an error wrapping multiple validation
// errors returned by DeleteJsonWebKeyRequest.ValidateAll() if the designated
// constraints aren't met.
type DeleteJsonWebKeyRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteJsonWebKeyRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteJsonWebKeyRequestMultiError) AllErrors() []error { return m }

// DeleteJsonWebKeyRequestValidationError is the validation error returned by
// DeleteJsonWebKeyRequest.Validate if the designated constraints aren't met.
type DeleteJsonWebKeyRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteJsonWebKeyRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteJsonWebKeyRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteJsonWebKeyRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteJsonWebKeyRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteJsonWebKeyRequestValidationError) ErrorName() string {
	return "DeleteJsonWebKeyRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteJsonWebKeyRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteJsonWebKeyRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteJsonWebKeyRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteJsonWebKeyRequestValidationError{}

// Validate checks the field values on GetJsonWebKeyRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetJsonWebKeyRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetJsonWebKeyRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetJsonWebKeyRequestMultiError, or nil if none found.
func (m *GetJsonWebKeyRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetJsonWebKeyRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetSet()) < 1 {
		err := GetJsonWebKeyRequestValidationError{
			field:  "Set",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetKid()) < 1 {
		err := GetJsonWebKeyRequestValidationError{
			field:  "Kid",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetJsonWebKeyRequestMultiError(errors)
	}

	return nil
}

// GetJsonWebKeyRequestMultiError is an error wrapping multiple validation
// errors returned by GetJsonWebKeyRequest.ValidateAll() if the designated
// constraints aren't met.
type GetJsonWebKeyRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetJsonWebKeyRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetJsonWebKeyRequestMultiError) AllErrors() []error { return m }

// GetJsonWebKeyRequestValidationError is the validation error returned by
// GetJsonWebKeyRequest.Validate if the designated constraints aren't met.
type GetJsonWebKeyRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetJsonWebKeyRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetJsonWebKeyRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetJsonWebKeyRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetJsonWebKeyRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetJsonWebKeyRequestValidationError) ErrorName() string {
	return "GetJsonWebKeyRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetJsonWebKeyRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetJsonWebKeyRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetJsonWebKeyRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetJsonWebKeyRequestValidationError{}

// Validate checks the field values on UpdateJsonWebKeyRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateJsonWebKeyRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateJsonWebKeyRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateJsonWebKeyRequestMultiError, or nil if none found.
func (m *UpdateJsonWebKeyRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateJsonWebKeyRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetSet()) < 1 {
		err := UpdateJsonWebKeyRequestValidationError{
			field:  "Set",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetKid()) < 1 {
		err := UpdateJsonWebKeyRequestValidationError{
			field:  "Kid",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetKey()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UpdateJsonWebKeyRequestValidationError{
					field:  "Key",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UpdateJsonWebKeyRequestValidationError{
					field:  "Key",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetKey()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdateJsonWebKeyRequestValidationError{
				field:  "Key",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return UpdateJsonWebKeyRequestMultiError(errors)
	}

	return nil
}

// UpdateJsonWebKeyRequestMultiError is an error wrapping multiple validation
// errors returned by UpdateJsonWebKeyRequest.ValidateAll() if the designated
// constraints aren't met.
type UpdateJsonWebKeyRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateJsonWebKeyRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateJsonWebKeyRequestMultiError) AllErrors() []error { return m }

// UpdateJsonWebKeyRequestValidationError is the validation error returned by
// UpdateJsonWebKeyRequest.Validate if the designated constraints aren't met.
type UpdateJsonWebKeyRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateJsonWebKeyRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateJsonWebKeyRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateJsonWebKeyRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateJsonWebKeyRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateJsonWebKeyRequestValidationError) ErrorName() string {
	return "UpdateJsonWebKeyRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateJsonWebKeyRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateJsonWebKeyRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateJsonWebKeyRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateJsonWebKeyRequestValidationError{}

// Validate checks the field values on JsonWebKey with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *JsonWebKey) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on JsonWebKey with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in JsonWebKeyMultiError, or
// nil if none found.
func (m *JsonWebKey) ValidateAll() error {
	return m.validate(true)
}

func (m *JsonWebKey) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetAlg()) < 1 {
		err := JsonWebKeyValidationError{
			field:  "Alg",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetKid()) < 1 {
		err := JsonWebKeyValidationError{
			field:  "Kid",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetKty()) < 1 {
		err := JsonWebKeyValidationError{
			field:  "Kty",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetUse()) < 1 {
		err := JsonWebKeyValidationError{
			field:  "Use",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.Crv != nil {
		// no validation rules for Crv
	}

	if m.D != nil {
		// no validation rules for D
	}

	if m.Dp != nil {
		// no validation rules for Dp
	}

	if m.Dq != nil {
		// no validation rules for Dq
	}

	if m.E != nil {
		// no validation rules for E
	}

	if m.K != nil {
		// no validation rules for K
	}

	if m.N != nil {
		// no validation rules for N
	}

	if m.P != nil {
		// no validation rules for P
	}

	if m.Q != nil {
		// no validation rules for Q
	}

	if m.Qi != nil {
		// no validation rules for Qi
	}

	if m.X != nil {
		// no validation rules for X
	}

	if m.Y != nil {
		// no validation rules for Y
	}

	if len(errors) > 0 {
		return JsonWebKeyMultiError(errors)
	}

	return nil
}

// JsonWebKeyMultiError is an error wrapping multiple validation errors
// returned by JsonWebKey.ValidateAll() if the designated constraints aren't met.
type JsonWebKeyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m JsonWebKeyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m JsonWebKeyMultiError) AllErrors() []error { return m }

// JsonWebKeyValidationError is the validation error returned by
// JsonWebKey.Validate if the designated constraints aren't met.
type JsonWebKeyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e JsonWebKeyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e JsonWebKeyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e JsonWebKeyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e JsonWebKeyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e JsonWebKeyValidationError) ErrorName() string { return "JsonWebKeyValidationError" }

// Error satisfies the builtin error interface
func (e JsonWebKeyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sJsonWebKey.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = JsonWebKeyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = JsonWebKeyValidationError{}

// Validate checks the field values on JsonWebKeySet with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *JsonWebKeySet) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on JsonWebKeySet with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in JsonWebKeySetMultiError, or
// nil if none found.
func (m *JsonWebKeySet) ValidateAll() error {
	return m.validate(true)
}

func (m *JsonWebKeySet) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetKeys() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, JsonWebKeySetValidationError{
						field:  fmt.Sprintf("Keys[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, JsonWebKeySetValidationError{
						field:  fmt.Sprintf("Keys[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return JsonWebKeySetValidationError{
					field:  fmt.Sprintf("Keys[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return JsonWebKeySetMultiError(errors)
	}

	return nil
}

// JsonWebKeySetMultiError is an error wrapping multiple validation errors
// returned by JsonWebKeySet.ValidateAll() if the designated constraints
// aren't met.
type JsonWebKeySetMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m JsonWebKeySetMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m JsonWebKeySetMultiError) AllErrors() []error { return m }

// JsonWebKeySetValidationError is the validation error returned by
// JsonWebKeySet.Validate if the designated constraints aren't met.
type JsonWebKeySetValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e JsonWebKeySetValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e JsonWebKeySetValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e JsonWebKeySetValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e JsonWebKeySetValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e JsonWebKeySetValidationError) ErrorName() string { return "JsonWebKeySetValidationError" }

// Error satisfies the builtin error interface
func (e JsonWebKeySetValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sJsonWebKeySet.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = JsonWebKeySetValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = JsonWebKeySetValidationError{}
