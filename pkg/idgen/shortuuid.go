package idgen

import (
	"context"
	"github.com/lithammer/shortuuid/v3"
)

type ShortUuid struct {
}

func (u *ShortUuid) Gen(ctx context.Context) (string, error) {
	return shortuuid.New(), nil
}
