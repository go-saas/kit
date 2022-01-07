package biz

import "strings"

type LookupNormalizer interface {
	// Name normalizer
	Name(name string) string
	// Email normalizer
	Email(email string) string
}

type lookupNormalizer struct {
}

func NewLookupNormalizer() LookupNormalizer {
	return lookupNormalizer{}
}
func (l lookupNormalizer) Name(name string) string {
	return strings.ToUpper(name)
}

func (l lookupNormalizer) Email(email string) string {
	return l.Name(email)
}
