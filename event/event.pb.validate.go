// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: event/event.proto

package event

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

// Validate checks the field values on Config with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Config) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Config with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ConfigMultiError, or nil if none found.
func (m *Config) ValidateAll() error {
	return m.validate(true)
}

func (m *Config) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Type

	// no validation rules for Addr

	// no validation rules for Topic

	// no validation rules for Group

	if all {
		switch v := interface{}(m.GetKafka()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ConfigValidationError{
					field:  "Kafka",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ConfigValidationError{
					field:  "Kafka",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetKafka()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ConfigValidationError{
				field:  "Kafka",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetPulsar()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ConfigValidationError{
					field:  "Pulsar",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ConfigValidationError{
					field:  "Pulsar",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPulsar()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ConfigValidationError{
				field:  "Pulsar",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetExtra()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ConfigValidationError{
					field:  "Extra",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ConfigValidationError{
					field:  "Extra",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetExtra()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ConfigValidationError{
				field:  "Extra",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return ConfigMultiError(errors)
	}

	return nil
}

// ConfigMultiError is an error wrapping multiple validation errors returned by
// Config.ValidateAll() if the designated constraints aren't met.
type ConfigMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ConfigMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ConfigMultiError) AllErrors() []error { return m }

// ConfigValidationError is the validation error returned by Config.Validate if
// the designated constraints aren't met.
type ConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ConfigValidationError) ErrorName() string { return "ConfigValidationError" }

// Error satisfies the builtin error interface
func (e ConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ConfigValidationError{}

// Validate checks the field values on Config_Kafka with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Config_Kafka) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Config_Kafka with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in Config_KafkaMultiError, or
// nil if none found.
func (m *Config_Kafka) ValidateAll() error {
	return m.validate(true)
}

func (m *Config_Kafka) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.Version != nil {

		if all {
			switch v := interface{}(m.GetVersion()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, Config_KafkaValidationError{
						field:  "Version",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, Config_KafkaValidationError{
						field:  "Version",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetVersion()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return Config_KafkaValidationError{
					field:  "Version",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return Config_KafkaMultiError(errors)
	}

	return nil
}

// Config_KafkaMultiError is an error wrapping multiple validation errors
// returned by Config_Kafka.ValidateAll() if the designated constraints aren't met.
type Config_KafkaMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m Config_KafkaMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m Config_KafkaMultiError) AllErrors() []error { return m }

// Config_KafkaValidationError is the validation error returned by
// Config_Kafka.Validate if the designated constraints aren't met.
type Config_KafkaValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e Config_KafkaValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e Config_KafkaValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e Config_KafkaValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e Config_KafkaValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e Config_KafkaValidationError) ErrorName() string { return "Config_KafkaValidationError" }

// Error satisfies the builtin error interface
func (e Config_KafkaValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sConfig_Kafka.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = Config_KafkaValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = Config_KafkaValidationError{}

// Validate checks the field values on Config_Pulsar with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Config_Pulsar) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Config_Pulsar with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in Config_PulsarMultiError, or
// nil if none found.
func (m *Config_Pulsar) ValidateAll() error {
	return m.validate(true)
}

func (m *Config_Pulsar) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.OperationTimeout != nil {

		if all {
			switch v := interface{}(m.GetOperationTimeout()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, Config_PulsarValidationError{
						field:  "OperationTimeout",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, Config_PulsarValidationError{
						field:  "OperationTimeout",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetOperationTimeout()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return Config_PulsarValidationError{
					field:  "OperationTimeout",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.ConnectionTimeout != nil {

		if all {
			switch v := interface{}(m.GetConnectionTimeout()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, Config_PulsarValidationError{
						field:  "ConnectionTimeout",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, Config_PulsarValidationError{
						field:  "ConnectionTimeout",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetConnectionTimeout()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return Config_PulsarValidationError{
					field:  "ConnectionTimeout",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return Config_PulsarMultiError(errors)
	}

	return nil
}

// Config_PulsarMultiError is an error wrapping multiple validation errors
// returned by Config_Pulsar.ValidateAll() if the designated constraints
// aren't met.
type Config_PulsarMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m Config_PulsarMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m Config_PulsarMultiError) AllErrors() []error { return m }

// Config_PulsarValidationError is the validation error returned by
// Config_Pulsar.Validate if the designated constraints aren't met.
type Config_PulsarValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e Config_PulsarValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e Config_PulsarValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e Config_PulsarValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e Config_PulsarValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e Config_PulsarValidationError) ErrorName() string { return "Config_PulsarValidationError" }

// Error satisfies the builtin error interface
func (e Config_PulsarValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sConfig_Pulsar.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = Config_PulsarValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = Config_PulsarValidationError{}
