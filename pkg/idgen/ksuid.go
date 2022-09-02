package idgen

import (
	"context"
	"github.com/segmentio/ksuid"
)

type Ksuid struct {
}

func (u *Ksuid) Gen(ctx context.Context) (string, error) {
	uid := ksuid.New()
	return uid.String(), nil
}
