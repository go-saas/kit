package plugins

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authn"
	"github.com/goxiaoy/go-saas-kit/pkg/authn/jwt"
	conf2 "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"net/http"

	pkgHTTP "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/plugin"
)

func init() {
	err := plugin.RegisterPlugin(&KitAuthn{})
	if err != nil {
		log.Fatalf("failed to register plugin kit_authn: %s", err)
	}
}

type KitAuthn struct {
}

type KitAuthConf struct {
	Security   *conf2.Security
	Services   *conf2.Services
	ClientName string
}

var tokenizer jwt.Tokenizer
var tokenManager api.TokenManager
var apiClient *conf2.Client
var apiOpt *api.Option

func Init(t jwt.Tokenizer, tmr api.TokenManager, clientName api.ClientName, services *conf2.Services, ao *api.Option) error {
	tokenizer = t
	tokenManager = tmr
	clientCfg, ok := services.Clients[string(clientName)]
	if !ok {
		return errors.New(fmt.Sprintf(" %v client not found", clientName))
	}
	apiClient = clientCfg

	apiOpt = ao
	return nil
}

func (p *KitAuthn) Name() string {
	return "kit_authn"
}

func (p *KitAuthn) ParseConf(in []byte) (interface{}, error) {
	conf := KitAuthConf{}
	err := json.Unmarshal(in, &conf)
	if err != nil {
		return nil, err
	}
	//init all
	_, _, err = initApp(conf.Services, conf.Security, api.ClientName(conf.ClientName))
	if err != nil {
		return nil, err
	}
	return conf, err
}

func (p *KitAuthn) Filter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {

	ctx := context.Background()

	var t = ""
	if auth := r.Header().Get(jwt.AuthorizationHeader); len(auth) > 0 {
		t = jwt.ExtractHeaderToken(auth)
	}
	if len(t) == 0 {
		t = r.Args().Get(jwt.AuthorizationQuery)
	}
	uid := ""
	clientId := ""

	if claims, err := jwt.ExtractAndValidate(tokenizer, t); err != nil {
		log.Fatalf("fail to extract and validate token %s", err)
	} else {
		if claims.Subject != "" {
			uid = claims.Subject
		} else {
			uid = claims.Uid
		}
		clientId = claims.ClientId
	}

	if len(clientId) > 0 {
		ctx = authn.NewClientContext(ctx, clientId)
	}
	ctx = authn.NewUserContext(ctx, authn.NewUserInfo(uid))

	//set auth token

	//use token mgr
	token, err := tokenManager.GetOrGenerateToken(ctx, &conf2.Client{
		ClientId:     apiClient.ClientId,
		ClientSecret: apiClient.ClientSecret,
	})
	if err != nil {
		log.Fatalf("%s", err)
		w.WriteHeader(500)
		return
	}
	//replace with internal token
	r.Header().Set(jwt.AuthorizationHeader, fmt.Sprintf("%s %s", jwt.BearerTokenType, token))

	//recover header
	for _, contributor := range apiOpt.Contributor {
		headers := contributor.CreateHeader(ctx)
		if headers != nil {
			for k, v := range headers {
				w.Header().Set(fmt.Sprintf("%s%s", api.PrefixOrDefault(""), k), v)
			}
		}
	}

	//TODO session auth

	//continue request
	return
}
