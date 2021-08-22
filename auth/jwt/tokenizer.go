package jwt

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type Claims struct {
	Uid string `json:"id,omitempty"`
	jwt.StandardClaims
}

type Tokenizer interface {
	Issue(userId string) (token string, err error)
	Parse(token string) (claims *Claims, err error)
}

type tokenizer struct {
	config *TokenizerConfig
}

type TokenizerConfig struct {
	ExpireDuration time.Duration
	Secret         string
}

func NewTokenizer(c *TokenizerConfig) Tokenizer {
	return tokenizer{config: c}
}

var _ Tokenizer = (*tokenizer)(nil)

func (t tokenizer) Issue(userId string) (token string, err error) {
	claims := Claims{
		userId,
		jwt.StandardClaims{
			Id:        userId,
			Subject:   userId,
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(t.config.ExpireDuration).Unix(),
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(t.config.Secret))
	return
}

func (t tokenizer) Parse(token string) (claims *Claims, err error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.config.Secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
