package api

import (
	"context"
	"github.com/goxiaoy/go-saas-kit/pkg/auth/jwt"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"time"
)

type TokenManager interface {
	GetOrGenerateToken(ctx context.Context, client *conf.Client) (token string, err error)
}

// InMemoryTokenManager TODO should use centralize authorization server to get jwt
type InMemoryTokenManager struct {
	token         string
	currentClient string
	tokenizer     jwt.Tokenizer
}

var _ TokenManager = (*InMemoryTokenManager)(nil)

func NewInMemoryTokenManager(tokenizer jwt.Tokenizer) *InMemoryTokenManager {
	return &InMemoryTokenManager{tokenizer: tokenizer}
}

func (i *InMemoryTokenManager) GetOrGenerateToken(ctx context.Context, client *conf.Client) (string, error) {
	if client.ClientId == i.currentClient && i.token != "" {
		return i.token, nil
	}
	//1000 year token....
	token, err := i.tokenizer.Issue(jwt.NewClientClaim(client.GetClientId()), time.Hour*876000)
	if err != nil {
		return token, err
	}
	i.token = token
	i.currentClient = client.GetClientId()
	return token, nil
}
