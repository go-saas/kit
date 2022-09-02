package idgen

import (
	"context"
	"github.com/google/uuid"
)

type Uuid struct {
}

func (u *Uuid) Gen(ctx context.Context) (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}