package api

import (
	"context"
	"fmt"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn/jwt"
	"github.com/go-saas/kit/pkg/conf"
)

const ServiceName = "dtmservice"

var ClientConf = &conf.Client{
	ClientId: ServiceName,
}

func MustAddBranchHeader(ctx context.Context, tokenMgr sapi.TokenManager) map[string]string {
	t, err := tokenMgr.GetOrGenerateToken(ctx, ClientConf)
	if err != nil {
		panic(err)
	}
	return map[string]string{
		jwt.AuthorizationHeader: fmt.Sprintf("%s %s", jwt.BearerTokenType, t),
	}
}
