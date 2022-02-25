package api

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
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
	l             *log.Helper
}

var _ TokenManager = (*InMemoryTokenManager)(nil)

func NewInMemoryTokenManager(tokenizer jwt.Tokenizer, logger log.Logger) *InMemoryTokenManager {
	return &InMemoryTokenManager{tokenizer: tokenizer, l: log.NewHelper(log.With(logger, "module", "InMemoryTokenManager"))}
}

func (i *InMemoryTokenManager) GetOrGenerateToken(ctx context.Context, client *conf.Client) (string, error) {
	if client.ClientId == i.currentClient && i.token != "" {
		return i.token, nil
	}
	//TODO 1000 year token....
	token, err := i.tokenizer.Issue(jwt.NewClientClaim(client.GetClientId()), time.Hour*876000)
	if err != nil {
		return token, err
	}
	i.l.Info("generate token for client: %s", client.GetClientId())
	i.token = token
	i.currentClient = client.GetClientId()
	return token, nil
}
