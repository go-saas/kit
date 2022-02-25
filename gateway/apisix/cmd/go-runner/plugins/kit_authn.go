package plugins

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	klog "github.com/go-kratos/kratos/v2/log"
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
}

var (
	tokenizer    jwt.Tokenizer
	tokenManager api.TokenManager
	apiClient    *conf2.Client
	apiOpt       *api.Option
)

func Init(t jwt.Tokenizer, tmr api.TokenManager, clientName api.ClientName, services *conf2.Services, logger klog.Logger) error {
	tokenizer = t
	tokenManager = tmr
	clientCfg, ok := services.Clients[string(clientName)]
	if !ok {
		return errors.New(fmt.Sprintf(" %v client not found", clientName))
	}
	apiClient = clientCfg
	apiOpt = api.NewOption("", true, api.NewUserContributor(logger), api.NewClientContributor(false, logger))
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
	if len(t) > 0 {
		if claims, err := jwt.ExtractAndValidate(tokenizer, t); err != nil {
			log.Errorf("fail to extract and validate token %s", err)
		} else {
			if claims.Subject != "" {
				uid = claims.Subject
			} else {
				uid = claims.Uid
			}
			clientId = claims.ClientId
		}
	}

	//TODO session auth

	log.Infof("resolve user: %s client: %s", uid, clientId)
	//keep previous client id
	ctx = authn.NewClientContext(ctx, clientId)
	ctx = authn.NewUserContext(ctx, authn.NewUserInfo(uid))

	//set auth token
	//use token mgr
	token, err := tokenManager.GetOrGenerateToken(ctx, &conf2.Client{
		ClientId:     apiClient.ClientId,
		ClientSecret: apiClient.ClientSecret,
	})
	if err != nil {
		log.Errorf("%s", err)
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
				nh := fmt.Sprintf("%s%s", api.PrefixOrDefault(""), k)
				log.Infof("set header: %s value: %s", nh, v)
				r.Header().Set(nh, v)
			}
		}
	}

	//continue request
	return
}
