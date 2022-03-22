package biz

import (
	v1 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"github.com/nyaruka/phonenumbers"
	"net/mail"
	"strings"
)

type LookupNormalizer interface {
	// Name normalizer
	Name(name string) (string, error)
	// Email normalizer
	Email(email string) (string, error)
	// Phone normalizer
	Phone(phone string) (string, error)
}

type lookupNormalizer struct {
}

func NewLookupNormalizer() LookupNormalizer {
	return lookupNormalizer{}
}
func (l lookupNormalizer) Name(name string) (string, error) {

	if _, err := l.Email(name); err == nil {
		return "", v1.ErrorInvalidUsername("")
	}
	if _, err := l.Phone(name); err == nil {
		return "", v1.ErrorInvalidUsername("")
	}
	return strings.ToUpper(name), nil
}

func (l lookupNormalizer) Email(email string) (string, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return "", v1.ErrorInvalidEmail("%s", err)
	}
	return strings.ToUpper(email), nil
}

func (l lookupNormalizer) Phone(phone string) (string, error) {
	num, err := phonenumbers.Parse(phone, "US")
	if err != nil {
		return "", v1.ErrorInvalidPhone("%s", err)
	}
	if ok := phonenumbers.IsValidNumber(num); !ok {
		return "", v1.ErrorInvalidPhone("")
	}
	formattedNum := phonenumbers.Format(num, phonenumbers.NATIONAL)
	return formattedNum, err
}
