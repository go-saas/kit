package service

import (
	pb "github.com/goxiaoy/go-saas-kit/user/api/auth/v1"
	"github.com/goxiaoy/go-saas-kit/user/private/biz"
)

func ConvertError(err error) error {
	if err == biz.ErrUserLocked {
		return pb.ErrorUserLocked("user locked")
	}
	if err == biz.ErrInvalidCredential {
		return pb.ErrorInvalidCredentials("invalid credentials")
	}
	return err
}
