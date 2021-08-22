package biz

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	ctx := context.Background()
	h := NewPasswordHasher()
	hash, err := h.HashPassword(ctx, nil, "password")
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(hash, "$argon2id$v=19$m=65536,t=1,p=2$"), hash)
}

func TestVerifyHashedPassword(t *testing.T) {
	ctx := context.Background()
	h := NewPasswordHasher()
	hash, _ := h.HashPassword(ctx, nil, "password")
	r := h.VerifyHashedPassword(ctx, nil, hash, "password")
	assert.Equal(t, PasswordVerificationSuccess, r)
	rf := h.VerifyHashedPassword(ctx, nil, hash, "password1")
	assert.Equal(t, PasswordVerificationFail, rf)

}
