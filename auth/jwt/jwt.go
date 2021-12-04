package jwt

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewTokenizer)

const (
	AuthorizationHeader = "Authorization"
	BearerTokenType     = "Bearer"
)
