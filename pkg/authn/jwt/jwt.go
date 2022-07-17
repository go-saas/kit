package jwt

import (
	kitdi "github.com/go-saas/kit/pkg/di"
)

var ProviderSet = kitdi.NewSet(NewTokenizer, NewTokenizerConfig)

const (
	AuthorizationHeader = "Authorization"
	BearerTokenType     = "Bearer"
	AuthorizationQuery  = "access_token"
)
