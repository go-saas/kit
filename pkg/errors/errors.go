package errors

import "github.com/go-kratos/kratos/v2/errors"

func UnRecoverableError(err error) bool {
	fr := errors.FromError(err)
	if fr.Code < 500 {
		return false
	}
	return true
}
