package jwt

import (
	"github.com/go-saas/kit/pkg/conf"
	"github.com/golang-jwt/jwt"
	"time"
)

type Claims struct {
	Uid      string `json:"id,omitempty"`
	ClientId string `json:"clientId,omitempty"`
	jwt.StandardClaims
}

func NewUserClaim(userId string) *Claims {
	return &Claims{
		Uid: userId,
		StandardClaims: jwt.StandardClaims{
			Id:      userId,
			Subject: userId,
		},
	}
}

func NewClientClaim(clientId string) *Claims {
	return &Claims{
		ClientId: clientId,
	}
}

type Tokenizer interface {
	Issue(claims *Claims, duration time.Duration) (token string, err error)
	Parse(token string) (claims *Claims, err error)
}

type tokenizer struct {
	config *TokenizerConfig
}

type TokenizerConfig struct {
	Issuer         string
	ExpireDuration time.Duration
	Secret         string
}

func NewTokenizerConfig(c *conf.Security) *TokenizerConfig {
	return &TokenizerConfig{
		Issuer:         c.Jwt.Issuer,
		ExpireDuration: c.Jwt.ExpireIn.AsDuration(),
		Secret:         c.Jwt.Secret,
	}
}

func NewTokenizer(c *TokenizerConfig) Tokenizer {
	return &tokenizer{config: c}
}

var _ Tokenizer = (*tokenizer)(nil)

func (t *tokenizer) Issue(claims *Claims, duration time.Duration) (token string, err error) {
	claims.StandardClaims.NotBefore = time.Now().Unix()
	claims.StandardClaims.IssuedAt = time.Now().Unix()
	if duration > 0 {
		claims.StandardClaims.ExpiresAt = time.Now().Add(duration).Unix()
	} else {
		claims.StandardClaims.ExpiresAt = time.Now().Add(t.config.ExpireDuration).Unix()
	}
	claims.Issuer = t.config.Issuer
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(t.config.Secret))
	return
}

func (t *tokenizer) Parse(token string) (claims *Claims, err error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.config.Secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid && claims.VerifyIssuer(t.config.Issuer, false) {
			return claims, nil
		}
	}
	return nil, err
}
