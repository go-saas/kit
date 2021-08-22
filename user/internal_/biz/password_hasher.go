package biz

import (
	"context"
	"github.com/alexedwards/argon2id"
)

type PasswordHasher interface {
	HashPassword(ctx context.Context, user *User, password string) (hash string, err error)
	VerifyHashedPassword(ctx context.Context, user *User, hashedPassword string, providedPassword string) PasswordVerificationResult
}

type PasswordVerificationResult int32

const (
	PasswordVerificationFail PasswordVerificationResult = iota
	PasswordVerificationSuccess
	PasswordVerificationSuccessRehashNeeded
)

type passwordHasher struct {
}

func NewPasswordHasher() PasswordHasher {
	return passwordHasher{}
}

func (p passwordHasher) HashPassword(ctx context.Context, user *User, password string) (hash string, err error) {
	// CreateHash returns a Argon2id hash of a plain-text password using the
	// provided algorithm parameters. The returned hash follows the format used
	// by the Argon2 reference C implementation and looks like this:
	// $argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
	hash, err = argon2id.CreateHash(password, argon2id.DefaultParams)
	return
}

func (p passwordHasher) VerifyHashedPassword(ctx context.Context, user *User, hashedPassword string, providedPassword string) PasswordVerificationResult {
	// ComparePasswordAndHash performs a constant-time comparison between a
	// plain-text password and Argon2id hash, using the parameters and salt
	// contained in the hash. It returns true if they match, otherwise it returns
	// false.
	match, err := argon2id.ComparePasswordAndHash(providedPassword, hashedPassword)
	if err != nil {
		return PasswordVerificationFail
	}
	if match {
		return PasswordVerificationSuccess
	}
	return PasswordVerificationFail
}
