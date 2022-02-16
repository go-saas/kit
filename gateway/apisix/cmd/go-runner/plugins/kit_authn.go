package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	"net/http"

	pkgHTTP "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/plugin"
)

func init() {
	err := plugin.RegisterPlugin(&KitAuthn{})
	if err != nil {
		log.Fatalf("failed to register plugin say: %s", err)
	}
}

type KitAuthn struct {
}

type KitAuthConf struct {
	HeaderPrefix string `json:"header_prefix"`
	TenantKey    string `json:"tenant_key"`
}

func (p *KitAuthn) Name() string {
	return "kit_authn"
}

func (p *KitAuthn) ParseConf(in []byte) (interface{}, error) {
	conf := KitAuthConf{}
	err := json.Unmarshal(in, &conf)
	return conf, err
}

func (p *KitAuthn) Filter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
	cfg := conf.(KitAuthConf)

	ctx := context.Background()

	//TODO resolve auth with jwt and session

	apiOpt := api.NewOption(api.Prefix(cfg.HeaderPrefix), true, api.NewUserContributor())
	//by pass jwt token
	if rawToken, ok := jwt.FromJWTContext(ctx); ok {
		//bypass raw token
		w.Header().Set(jwt.AuthorizationHeader, fmt.Sprintf("%s %s", jwt.BearerTokenType, rawToken))
	}
	//recover header
	for _, contributor := range apiOpt.Contributor {
		headers := contributor.CreateHeader(ctx)
		if headers != nil {
			for k, v := range headers {
				w.Header().Set(fmt.Sprintf("%s%s", cfg.HeaderPrefix, k), v)
			}
		}
	}
	w.Header().Add("X-Resp-A6-Runner", "Go")
	//continue request
	return
}
