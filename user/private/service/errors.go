package service

import (
	"errors"
	pb "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	v1 "github.com/goxiaoy/go-saas-kit/user/api/user/v1"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
)

func ConvertError(err error) error {
	if err == biz.ErrUserLocked {
		return pb.ErrorUserLocked("user locked")
	}
	if err == biz.ErrInvalidCredential {
		return pb.ErrorInvalidCredentials("invalid credentials")
	}
	if errors.Is(err, biz.ErrInsufficientStrength) {
		return v1.ErrorPasswordInsufficientStrength("")
	}
	if errors.Is(err, biz.ErrDuplicateEmail) {
		return v1.ErrorDuplicateEmail("")
	}
	if errors.Is(err, biz.ErrDuplicateUsername) {
		return v1.ErrorDuplicateUsername("")
	}
	if errors.Is(err, biz.ErrDuplicatePhone) {
		return v1.ErrorDuplicatePhone("")
	}
	return err
}
