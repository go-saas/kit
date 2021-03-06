// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: user/private/biz/job.proto

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

// Validate checks the field values on UserMigrationTask with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *UserMigrationTask) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserMigrationTask with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UserMigrationTaskMultiError, or nil if none found.
func (m *UserMigrationTask) ValidateAll() error {
	return m.validate(true)
}

func (m *UserMigrationTask) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for AdminEmail

	// no validation rules for AdminUsername

	// no validation rules for AdminPassword

	// no validation rules for AdminUserId

	if len(errors) > 0 {
		return UserMigrationTaskMultiError(errors)
	}

	return nil
}

// UserMigrationTaskMultiError is an error wrapping multiple validation errors
// returned by UserMigrationTask.ValidateAll() if the designated constraints
// aren't met.
type UserMigrationTaskMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserMigrationTaskMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserMigrationTaskMultiError) AllErrors() []error { return m }

// UserMigrationTaskValidationError is the validation error returned by
// UserMigrationTask.Validate if the designated constraints aren't met.
type UserMigrationTaskValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserMigrationTaskValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserMigrationTaskValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserMigrationTaskValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserMigrationTaskValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserMigrationTaskValidationError) ErrorName() string {
	return "UserMigrationTaskValidationError"
}

// Error satisfies the builtin error interface
func (e UserMigrationTaskValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserMigrationTask.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserMigrationTaskValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserMigrationTaskValidationError{}
