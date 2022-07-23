// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: oidc/api/client/v1/client.proto

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

// Validate checks the field values on ListClientRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ListClientRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListClientRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListClientRequestMultiError, or nil if none found.
func (m *ListClientRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListClientRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Limit

	// no validation rules for Offset

	// no validation rules for ClientName

	// no validation rules for Owner

	if len(errors) > 0 {
		return ListClientRequestMultiError(errors)
	}

	return nil
}

// ListClientRequestMultiError is an error wrapping multiple validation errors
// returned by ListClientRequest.ValidateAll() if the designated constraints
// aren't met.
type ListClientRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListClientRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListClientRequestMultiError) AllErrors() []error { return m }

// ListClientRequestValidationError is the validation error returned by
// ListClientRequest.Validate if the designated constraints aren't met.
type ListClientRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListClientRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListClientRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListClientRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListClientRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListClientRequestValidationError) ErrorName() string {
	return "ListClientRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListClientRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListClientRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListClientRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListClientRequestValidationError{}

// Validate checks the field values on OAuth2ClientList with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *OAuth2ClientList) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OAuth2ClientList with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// OAuth2ClientListMultiError, or nil if none found.
func (m *OAuth2ClientList) ValidateAll() error {
	return m.validate(true)
}

func (m *OAuth2ClientList) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetItems() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, OAuth2ClientListValidationError{
						field:  fmt.Sprintf("Items[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, OAuth2ClientListValidationError{
						field:  fmt.Sprintf("Items[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return OAuth2ClientListValidationError{
					field:  fmt.Sprintf("Items[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return OAuth2ClientListMultiError(errors)
	}

	return nil
}

// OAuth2ClientListMultiError is an error wrapping multiple validation errors
// returned by OAuth2ClientList.ValidateAll() if the designated constraints
// aren't met.
type OAuth2ClientListMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OAuth2ClientListMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OAuth2ClientListMultiError) AllErrors() []error { return m }

// OAuth2ClientListValidationError is the validation error returned by
// OAuth2ClientList.Validate if the designated constraints aren't met.
type OAuth2ClientListValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OAuth2ClientListValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OAuth2ClientListValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OAuth2ClientListValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OAuth2ClientListValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OAuth2ClientListValidationError) ErrorName() string { return "OAuth2ClientListValidationError" }

// Error satisfies the builtin error interface
func (e OAuth2ClientListValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOAuth2ClientList.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OAuth2ClientListValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OAuth2ClientListValidationError{}

// Validate checks the field values on OAuth2Client with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *OAuth2Client) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OAuth2Client with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in OAuth2ClientMultiError, or
// nil if none found.
func (m *OAuth2Client) ValidateAll() error {
	return m.validate(true)
}

func (m *OAuth2Client) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, OAuth2ClientValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, OAuth2ClientValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return OAuth2ClientValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetUpdatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, OAuth2ClientValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, OAuth2ClientValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return OAuth2ClientValidationError{
				field:  "UpdatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.BackchannelLogoutSessionRequired != nil {
		// no validation rules for BackchannelLogoutSessionRequired
	}

	if m.BackchannelLogoutUri != nil {
		// no validation rules for BackchannelLogoutUri
	}

	if m.ClientId != nil {
		// no validation rules for ClientId
	}

	if m.ClientName != nil {
		// no validation rules for ClientName
	}

	if m.ClientSecret != nil {
		// no validation rules for ClientSecret
	}

	if m.ClientSecretExpiresAt != nil {
		// no validation rules for ClientSecretExpiresAt
	}

	if m.ClientUri != nil {
		// no validation rules for ClientUri
	}

	if m.FrontchannelLogoutSessionRequired != nil {
		// no validation rules for FrontchannelLogoutSessionRequired
	}

	if m.FrontchannelLogoutUri != nil {
		// no validation rules for FrontchannelLogoutUri
	}

	if m.Jwks != nil {

		if all {
			switch v := interface{}(m.GetJwks()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, OAuth2ClientValidationError{
						field:  "Jwks",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, OAuth2ClientValidationError{
						field:  "Jwks",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetJwks()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return OAuth2ClientValidationError{
					field:  "Jwks",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.JwksUri != nil {
		// no validation rules for JwksUri
	}

	if m.LogoUri != nil {
		// no validation rules for LogoUri
	}

	if m.Metadata != nil {

		if all {
			switch v := interface{}(m.GetMetadata()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, OAuth2ClientValidationError{
						field:  "Metadata",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, OAuth2ClientValidationError{
						field:  "Metadata",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetMetadata()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return OAuth2ClientValidationError{
					field:  "Metadata",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Owner != nil {
		// no validation rules for Owner
	}

	if m.PolicyUri != nil {
		// no validation rules for PolicyUri
	}

	if m.RegistrationAccessToken != nil {
		// no validation rules for RegistrationAccessToken
	}

	if m.RegistrationClientUri != nil {
		// no validation rules for RegistrationClientUri
	}

	if m.RequestObjectSigningAlg != nil {
		// no validation rules for RequestObjectSigningAlg
	}

	if m.Scope != nil {
		// no validation rules for Scope
	}

	if m.SectorIdentifierUri != nil {
		// no validation rules for SectorIdentifierUri
	}

	if m.SubjectType != nil {
		// no validation rules for SubjectType
	}

	if m.TokenEndpointAuthMethod != nil {
		// no validation rules for TokenEndpointAuthMethod
	}

	if m.TokenEndpointAuthSigningAlg != nil {
		// no validation rules for TokenEndpointAuthSigningAlg
	}

	if m.TosUri != nil {
		// no validation rules for TosUri
	}

	if m.UserinfoSignedResponseAlg != nil {
		// no validation rules for UserinfoSignedResponseAlg
	}

	if len(errors) > 0 {
		return OAuth2ClientMultiError(errors)
	}

	return nil
}

// OAuth2ClientMultiError is an error wrapping multiple validation errors
// returned by OAuth2Client.ValidateAll() if the designated constraints aren't met.
type OAuth2ClientMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OAuth2ClientMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OAuth2ClientMultiError) AllErrors() []error { return m }

// OAuth2ClientValidationError is the validation error returned by
// OAuth2Client.Validate if the designated constraints aren't met.
type OAuth2ClientValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OAuth2ClientValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OAuth2ClientValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OAuth2ClientValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OAuth2ClientValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OAuth2ClientValidationError) ErrorName() string { return "OAuth2ClientValidationError" }

// Error satisfies the builtin error interface
func (e OAuth2ClientValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOAuth2Client.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OAuth2ClientValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OAuth2ClientValidationError{}

// Validate checks the field values on DeleteOAuth2ClientRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteOAuth2ClientRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteOAuth2ClientRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteOAuth2ClientRequestMultiError, or nil if none found.
func (m *DeleteOAuth2ClientRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteOAuth2ClientRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := DeleteOAuth2ClientRequestValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DeleteOAuth2ClientRequestMultiError(errors)
	}

	return nil
}

// DeleteOAuth2ClientRequestMultiError is an error wrapping multiple validation
// errors returned by DeleteOAuth2ClientRequest.ValidateAll() if the
// designated constraints aren't met.
type DeleteOAuth2ClientRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteOAuth2ClientRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteOAuth2ClientRequestMultiError) AllErrors() []error { return m }

// DeleteOAuth2ClientRequestValidationError is the validation error returned by
// DeleteOAuth2ClientRequest.Validate if the designated constraints aren't met.
type DeleteOAuth2ClientRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteOAuth2ClientRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteOAuth2ClientRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteOAuth2ClientRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteOAuth2ClientRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteOAuth2ClientRequestValidationError) ErrorName() string {
	return "DeleteOAuth2ClientRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteOAuth2ClientRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteOAuth2ClientRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteOAuth2ClientRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteOAuth2ClientRequestValidationError{}

// Validate checks the field values on GetOAuth2ClientRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetOAuth2ClientRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetOAuth2ClientRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetOAuth2ClientRequestMultiError, or nil if none found.
func (m *GetOAuth2ClientRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetOAuth2ClientRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := GetOAuth2ClientRequestValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetOAuth2ClientRequestMultiError(errors)
	}

	return nil
}

// GetOAuth2ClientRequestMultiError is an error wrapping multiple validation
// errors returned by GetOAuth2ClientRequest.ValidateAll() if the designated
// constraints aren't met.
type GetOAuth2ClientRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetOAuth2ClientRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetOAuth2ClientRequestMultiError) AllErrors() []error { return m }

// GetOAuth2ClientRequestValidationError is the validation error returned by
// GetOAuth2ClientRequest.Validate if the designated constraints aren't met.
type GetOAuth2ClientRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetOAuth2ClientRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetOAuth2ClientRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetOAuth2ClientRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetOAuth2ClientRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetOAuth2ClientRequestValidationError) ErrorName() string {
	return "GetOAuth2ClientRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetOAuth2ClientRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetOAuth2ClientRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetOAuth2ClientRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetOAuth2ClientRequestValidationError{}

// Validate checks the field values on PatchOAuth2ClientRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *PatchOAuth2ClientRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PatchOAuth2ClientRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PatchOAuth2ClientRequestMultiError, or nil if none found.
func (m *PatchOAuth2ClientRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *PatchOAuth2ClientRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := PatchOAuth2ClientRequestValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetClient() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PatchOAuth2ClientRequestValidationError{
						field:  fmt.Sprintf("Client[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PatchOAuth2ClientRequestValidationError{
						field:  fmt.Sprintf("Client[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PatchOAuth2ClientRequestValidationError{
					field:  fmt.Sprintf("Client[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return PatchOAuth2ClientRequestMultiError(errors)
	}

	return nil
}

// PatchOAuth2ClientRequestMultiError is an error wrapping multiple validation
// errors returned by PatchOAuth2ClientRequest.ValidateAll() if the designated
// constraints aren't met.
type PatchOAuth2ClientRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PatchOAuth2ClientRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PatchOAuth2ClientRequestMultiError) AllErrors() []error { return m }

// PatchOAuth2ClientRequestValidationError is the validation error returned by
// PatchOAuth2ClientRequest.Validate if the designated constraints aren't met.
type PatchOAuth2ClientRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PatchOAuth2ClientRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PatchOAuth2ClientRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PatchOAuth2ClientRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PatchOAuth2ClientRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PatchOAuth2ClientRequestValidationError) ErrorName() string {
	return "PatchOAuth2ClientRequestValidationError"
}

// Error satisfies the builtin error interface
func (e PatchOAuth2ClientRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPatchOAuth2ClientRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PatchOAuth2ClientRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PatchOAuth2ClientRequestValidationError{}

// Validate checks the field values on PatchOAuth2Client with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *PatchOAuth2Client) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PatchOAuth2Client with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PatchOAuth2ClientMultiError, or nil if none found.
func (m *PatchOAuth2Client) ValidateAll() error {
	return m.validate(true)
}

func (m *PatchOAuth2Client) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Op

	// no validation rules for Path

	if all {
		switch v := interface{}(m.GetValue()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PatchOAuth2ClientValidationError{
					field:  "Value",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PatchOAuth2ClientValidationError{
					field:  "Value",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetValue()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PatchOAuth2ClientValidationError{
				field:  "Value",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.From != nil {
		// no validation rules for From
	}

	if len(errors) > 0 {
		return PatchOAuth2ClientMultiError(errors)
	}

	return nil
}

// PatchOAuth2ClientMultiError is an error wrapping multiple validation errors
// returned by PatchOAuth2Client.ValidateAll() if the designated constraints
// aren't met.
type PatchOAuth2ClientMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PatchOAuth2ClientMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PatchOAuth2ClientMultiError) AllErrors() []error { return m }

// PatchOAuth2ClientValidationError is the validation error returned by
// PatchOAuth2Client.Validate if the designated constraints aren't met.
type PatchOAuth2ClientValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PatchOAuth2ClientValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PatchOAuth2ClientValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PatchOAuth2ClientValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PatchOAuth2ClientValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PatchOAuth2ClientValidationError) ErrorName() string {
	return "PatchOAuth2ClientValidationError"
}

// Error satisfies the builtin error interface
func (e PatchOAuth2ClientValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPatchOAuth2Client.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PatchOAuth2ClientValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PatchOAuth2ClientValidationError{}

// Validate checks the field values on UpdateOAuth2ClientRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateOAuth2ClientRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateOAuth2ClientRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateOAuth2ClientRequestMultiError, or nil if none found.
func (m *UpdateOAuth2ClientRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateOAuth2ClientRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := UpdateOAuth2ClientRequestValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetClient()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UpdateOAuth2ClientRequestValidationError{
					field:  "Client",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UpdateOAuth2ClientRequestValidationError{
					field:  "Client",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetClient()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdateOAuth2ClientRequestValidationError{
				field:  "Client",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return UpdateOAuth2ClientRequestMultiError(errors)
	}

	return nil
}

// UpdateOAuth2ClientRequestMultiError is an error wrapping multiple validation
// errors returned by UpdateOAuth2ClientRequest.ValidateAll() if the
// designated constraints aren't met.
type UpdateOAuth2ClientRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateOAuth2ClientRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateOAuth2ClientRequestMultiError) AllErrors() []error { return m }

// UpdateOAuth2ClientRequestValidationError is the validation error returned by
// UpdateOAuth2ClientRequest.Validate if the designated constraints aren't met.
type UpdateOAuth2ClientRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateOAuth2ClientRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateOAuth2ClientRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateOAuth2ClientRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateOAuth2ClientRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateOAuth2ClientRequestValidationError) ErrorName() string {
	return "UpdateOAuth2ClientRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateOAuth2ClientRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateOAuth2ClientRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateOAuth2ClientRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateOAuth2ClientRequestValidationError{}