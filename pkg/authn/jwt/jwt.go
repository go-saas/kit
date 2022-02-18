package jwt

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewTokenizer, NewTokenizerConfig)

const (
	AuthorizationHeader = "Authorization"
	BearerTokenType     = "Bearer"
	AuthorizationQuery  = "access_token"
)
