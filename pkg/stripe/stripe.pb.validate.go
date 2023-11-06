// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: stripe/stripe.proto

package stripe

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

// Validate checks the field values on Conf with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Conf) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Conf with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ConfMultiError, or nil if none found.
func (m *Conf) ValidateAll() error {
	return m.validate(true)
}

func (m *Conf) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for IsTest

	// no validation rules for PublishKey

	// no validation rules for PrivateKey

	// no validation rules for WebhookKey

	// no validation rules for PriceTables

	if len(errors) > 0 {
		return ConfMultiError(errors)
	}

	return nil
}

// ConfMultiError is an error wrapping multiple validation errors returned by
// Conf.ValidateAll() if the designated constraints aren't met.
type ConfMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ConfMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ConfMultiError) AllErrors() []error { return m }

// ConfValidationError is the validation error returned by Conf.Validate if the
// designated constraints aren't met.
type ConfValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ConfValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ConfValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ConfValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ConfValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ConfValidationError) ErrorName() string { return "ConfValidationError" }

// Error satisfies the builtin error interface
func (e ConfValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sConf.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ConfValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ConfValidationError{}

// Validate checks the field values on Invoice with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Invoice) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Invoice with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in InvoiceMultiError, or nil if none found.
func (m *Invoice) ValidateAll() error {
	return m.validate(true)
}

func (m *Invoice) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if all {
		switch v := interface{}(m.GetPaymentIntent()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, InvoiceValidationError{
					field:  "PaymentIntent",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, InvoiceValidationError{
					field:  "PaymentIntent",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPaymentIntent()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return InvoiceValidationError{
				field:  "PaymentIntent",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return InvoiceMultiError(errors)
	}

	return nil
}

// InvoiceMultiError is an error wrapping multiple validation errors returned
// by Invoice.ValidateAll() if the designated constraints aren't met.
type InvoiceMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InvoiceMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InvoiceMultiError) AllErrors() []error { return m }

// InvoiceValidationError is the validation error returned by Invoice.Validate
// if the designated constraints aren't met.
type InvoiceValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InvoiceValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InvoiceValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InvoiceValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InvoiceValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InvoiceValidationError) ErrorName() string { return "InvoiceValidationError" }

// Error satisfies the builtin error interface
func (e InvoiceValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInvoice.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InvoiceValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InvoiceValidationError{}

// Validate checks the field values on PaymentIntent with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PaymentIntent) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PaymentIntent with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PaymentIntentMultiError, or
// nil if none found.
func (m *PaymentIntent) ValidateAll() error {
	return m.validate(true)
}

func (m *PaymentIntent) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for ClientSecret

	// no validation rules for Status

	if len(errors) > 0 {
		return PaymentIntentMultiError(errors)
	}

	return nil
}

// PaymentIntentMultiError is an error wrapping multiple validation errors
// returned by PaymentIntent.ValidateAll() if the designated constraints
// aren't met.
type PaymentIntentMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PaymentIntentMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PaymentIntentMultiError) AllErrors() []error { return m }

// PaymentIntentValidationError is the validation error returned by
// PaymentIntent.Validate if the designated constraints aren't met.
type PaymentIntentValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PaymentIntentValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PaymentIntentValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PaymentIntentValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PaymentIntentValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PaymentIntentValidationError) ErrorName() string { return "PaymentIntentValidationError" }

// Error satisfies the builtin error interface
func (e PaymentIntentValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPaymentIntent.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PaymentIntentValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PaymentIntentValidationError{}

// Validate checks the field values on Subscription with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Subscription) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Subscription with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SubscriptionMultiError, or
// nil if none found.
func (m *Subscription) ValidateAll() error {
	return m.validate(true)
}

func (m *Subscription) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if all {
		switch v := interface{}(m.GetLatestInvoice()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SubscriptionValidationError{
					field:  "LatestInvoice",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SubscriptionValidationError{
					field:  "LatestInvoice",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetLatestInvoice()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SubscriptionValidationError{
				field:  "LatestInvoice",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return SubscriptionMultiError(errors)
	}

	return nil
}

// SubscriptionMultiError is an error wrapping multiple validation errors
// returned by Subscription.ValidateAll() if the designated constraints aren't met.
type SubscriptionMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SubscriptionMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SubscriptionMultiError) AllErrors() []error { return m }

// SubscriptionValidationError is the validation error returned by
// Subscription.Validate if the designated constraints aren't met.
type SubscriptionValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SubscriptionValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SubscriptionValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SubscriptionValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SubscriptionValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SubscriptionValidationError) ErrorName() string { return "SubscriptionValidationError" }

// Error satisfies the builtin error interface
func (e SubscriptionValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSubscription.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SubscriptionValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SubscriptionValidationError{}

// Validate checks the field values on EphemeralKey with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *EphemeralKey) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on EphemeralKey with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in EphemeralKeyMultiError, or
// nil if none found.
func (m *EphemeralKey) ValidateAll() error {
	return m.validate(true)
}

func (m *EphemeralKey) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Secret

	if len(errors) > 0 {
		return EphemeralKeyMultiError(errors)
	}

	return nil
}

// EphemeralKeyMultiError is an error wrapping multiple validation errors
// returned by EphemeralKey.ValidateAll() if the designated constraints aren't met.
type EphemeralKeyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EphemeralKeyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EphemeralKeyMultiError) AllErrors() []error { return m }

// EphemeralKeyValidationError is the validation error returned by
// EphemeralKey.Validate if the designated constraints aren't met.
type EphemeralKeyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EphemeralKeyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EphemeralKeyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EphemeralKeyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EphemeralKeyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EphemeralKeyValidationError) ErrorName() string { return "EphemeralKeyValidationError" }

// Error satisfies the builtin error interface
func (e EphemeralKeyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEphemeralKey.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EphemeralKeyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EphemeralKeyValidationError{}
