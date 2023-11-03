// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: product/api/product/v1/product_internal.proto

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

// Validate checks the field values on CreateInternalProductRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateInternalProductRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateInternalProductRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateInternalProductRequestMultiError, or nil if none found.
func (m *CreateInternalProductRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateInternalProductRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetTitle()) < 1 {
		err := CreateInternalProductRequestValidationError{
			field:  "Title",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for ShortDesc

	// no validation rules for Desc

	if all {
		switch v := interface{}(m.GetMainPic()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "MainPic",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "MainPic",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMainPic()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateInternalProductRequestValidationError{
				field:  "MainPic",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetMedias() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("Medias[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("Medias[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CreateInternalProductRequestValidationError{
					field:  fmt.Sprintf("Medias[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetBadges() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("Badges[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("Badges[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CreateInternalProductRequestValidationError{
					field:  fmt.Sprintf("Badges[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if all {
		switch v := interface{}(m.GetVisibleFrom()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "VisibleFrom",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "VisibleFrom",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetVisibleFrom()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateInternalProductRequestValidationError{
				field:  "VisibleFrom",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetVisibleTo()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "VisibleTo",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "VisibleTo",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetVisibleTo()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateInternalProductRequestValidationError{
				field:  "VisibleTo",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for IsNew

	// no validation rules for Barcode

	// no validation rules for Model

	// no validation rules for IsSaleable

	if all {
		switch v := interface{}(m.GetSaleableFrom()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "SaleableFrom",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "SaleableFrom",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetSaleableFrom()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateInternalProductRequestValidationError{
				field:  "SaleableFrom",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetSaleableTo()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "SaleableTo",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "SaleableTo",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetSaleableTo()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateInternalProductRequestValidationError{
				field:  "SaleableTo",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetKeywords() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("Keywords[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("Keywords[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CreateInternalProductRequestValidationError{
					field:  fmt.Sprintf("Keywords[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetPrices() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("Prices[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("Prices[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CreateInternalProductRequestValidationError{
					field:  fmt.Sprintf("Prices[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for IsGiveaway

	for idx, item := range m.GetAttributes() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("Attributes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("Attributes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CreateInternalProductRequestValidationError{
					field:  fmt.Sprintf("Attributes[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for MultiSku

	for idx, item := range m.GetCampaignRules() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("CampaignRules[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("CampaignRules[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CreateInternalProductRequestValidationError{
					field:  fmt.Sprintf("CampaignRules[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for NeedShipping

	for idx, item := range m.GetStocks() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("Stocks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CreateInternalProductRequestValidationError{
						field:  fmt.Sprintf("Stocks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CreateInternalProductRequestValidationError{
					field:  fmt.Sprintf("Stocks[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if all {
		switch v := interface{}(m.GetContent()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "Content",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "Content",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetContent()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateInternalProductRequestValidationError{
				field:  "Content",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetManageInfo()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "ManageInfo",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreateInternalProductRequestValidationError{
					field:  "ManageInfo",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetManageInfo()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateInternalProductRequestValidationError{
				field:  "ManageInfo",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Active

	if m.MainCategoryKey != nil {
		// no validation rules for MainCategoryKey
	}

	if m.BrandId != nil {
		// no validation rules for BrandId
	}

	if len(errors) > 0 {
		return CreateInternalProductRequestMultiError(errors)
	}

	return nil
}

// CreateInternalProductRequestMultiError is an error wrapping multiple
// validation errors returned by CreateInternalProductRequest.ValidateAll() if
// the designated constraints aren't met.
type CreateInternalProductRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateInternalProductRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateInternalProductRequestMultiError) AllErrors() []error { return m }

// CreateInternalProductRequestValidationError is the validation error returned
// by CreateInternalProductRequest.Validate if the designated constraints
// aren't met.
type CreateInternalProductRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateInternalProductRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateInternalProductRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateInternalProductRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateInternalProductRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateInternalProductRequestValidationError) ErrorName() string {
	return "CreateInternalProductRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateInternalProductRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateInternalProductRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateInternalProductRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateInternalProductRequestValidationError{}

// Validate checks the field values on UpdateInternalProductRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateInternalProductRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateInternalProductRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateInternalProductRequestMultiError, or nil if none found.
func (m *UpdateInternalProductRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateInternalProductRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetProduct() == nil {
		err := UpdateInternalProductRequestValidationError{
			field:  "Product",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetProduct()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UpdateInternalProductRequestValidationError{
					field:  "Product",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UpdateInternalProductRequestValidationError{
					field:  "Product",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetProduct()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdateInternalProductRequestValidationError{
				field:  "Product",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetUpdateMask()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UpdateInternalProductRequestValidationError{
					field:  "UpdateMask",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UpdateInternalProductRequestValidationError{
					field:  "UpdateMask",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdateMask()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdateInternalProductRequestValidationError{
				field:  "UpdateMask",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return UpdateInternalProductRequestMultiError(errors)
	}

	return nil
}

// UpdateInternalProductRequestMultiError is an error wrapping multiple
// validation errors returned by UpdateInternalProductRequest.ValidateAll() if
// the designated constraints aren't met.
type UpdateInternalProductRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateInternalProductRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateInternalProductRequestMultiError) AllErrors() []error { return m }

// UpdateInternalProductRequestValidationError is the validation error returned
// by UpdateInternalProductRequest.Validate if the designated constraints
// aren't met.
type UpdateInternalProductRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateInternalProductRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateInternalProductRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateInternalProductRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateInternalProductRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateInternalProductRequestValidationError) ErrorName() string {
	return "UpdateInternalProductRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateInternalProductRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateInternalProductRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateInternalProductRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateInternalProductRequestValidationError{}

// Validate checks the field values on GetInternalProductRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetInternalProductRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetInternalProductRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetInternalProductRequestMultiError, or nil if none found.
func (m *GetInternalProductRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetInternalProductRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := GetInternalProductRequestValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetInternalProductRequestMultiError(errors)
	}

	return nil
}

// GetInternalProductRequestMultiError is an error wrapping multiple validation
// errors returned by GetInternalProductRequest.ValidateAll() if the
// designated constraints aren't met.
type GetInternalProductRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetInternalProductRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetInternalProductRequestMultiError) AllErrors() []error { return m }

// GetInternalProductRequestValidationError is the validation error returned by
// GetInternalProductRequest.Validate if the designated constraints aren't met.
type GetInternalProductRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetInternalProductRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetInternalProductRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetInternalProductRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetInternalProductRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetInternalProductRequestValidationError) ErrorName() string {
	return "GetInternalProductRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetInternalProductRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetInternalProductRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetInternalProductRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetInternalProductRequestValidationError{}

// Validate checks the field values on DeleteInternalProductRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteInternalProductRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteInternalProductRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteInternalProductRequestMultiError, or nil if none found.
func (m *DeleteInternalProductRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteInternalProductRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if len(errors) > 0 {
		return DeleteInternalProductRequestMultiError(errors)
	}

	return nil
}

// DeleteInternalProductRequestMultiError is an error wrapping multiple
// validation errors returned by DeleteInternalProductRequest.ValidateAll() if
// the designated constraints aren't met.
type DeleteInternalProductRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteInternalProductRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteInternalProductRequestMultiError) AllErrors() []error { return m }

// DeleteInternalProductRequestValidationError is the validation error returned
// by DeleteInternalProductRequest.Validate if the designated constraints
// aren't met.
type DeleteInternalProductRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteInternalProductRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteInternalProductRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteInternalProductRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteInternalProductRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteInternalProductRequestValidationError) ErrorName() string {
	return "DeleteInternalProductRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteInternalProductRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteInternalProductRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteInternalProductRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteInternalProductRequestValidationError{}

// Validate checks the field values on DeleteInternalProductReply with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteInternalProductReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteInternalProductReply with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteInternalProductReplyMultiError, or nil if none found.
func (m *DeleteInternalProductReply) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteInternalProductReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if len(errors) > 0 {
		return DeleteInternalProductReplyMultiError(errors)
	}

	return nil
}

// DeleteInternalProductReplyMultiError is an error wrapping multiple
// validation errors returned by DeleteInternalProductReply.ValidateAll() if
// the designated constraints aren't met.
type DeleteInternalProductReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteInternalProductReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteInternalProductReplyMultiError) AllErrors() []error { return m }

// DeleteInternalProductReplyValidationError is the validation error returned
// by DeleteInternalProductReply.Validate if the designated constraints aren't met.
type DeleteInternalProductReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteInternalProductReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteInternalProductReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteInternalProductReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteInternalProductReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteInternalProductReplyValidationError) ErrorName() string {
	return "DeleteInternalProductReplyValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteInternalProductReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteInternalProductReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteInternalProductReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteInternalProductReplyValidationError{}

// Validate checks the field values on GetInternalPriceRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetInternalPriceRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetInternalPriceRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetInternalPriceRequestMultiError, or nil if none found.
func (m *GetInternalPriceRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetInternalPriceRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := GetInternalPriceRequestValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetInternalPriceRequestMultiError(errors)
	}

	return nil
}

// GetInternalPriceRequestMultiError is an error wrapping multiple validation
// errors returned by GetInternalPriceRequest.ValidateAll() if the designated
// constraints aren't met.
type GetInternalPriceRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetInternalPriceRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetInternalPriceRequestMultiError) AllErrors() []error { return m }

// GetInternalPriceRequestValidationError is the validation error returned by
// GetInternalPriceRequest.Validate if the designated constraints aren't met.
type GetInternalPriceRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetInternalPriceRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetInternalPriceRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetInternalPriceRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetInternalPriceRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetInternalPriceRequestValidationError) ErrorName() string {
	return "GetInternalPriceRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetInternalPriceRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetInternalPriceRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetInternalPriceRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetInternalPriceRequestValidationError{}
