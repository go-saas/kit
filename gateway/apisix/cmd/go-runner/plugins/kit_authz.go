package plugins

import (
	"encoding/json"
	"fmt"
	pkgHTTP "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/plugin"
	"github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/samber/lo"
	"net/http"
	"strings"
)

// KitAuthz authorization plugin
type KitAuthz struct {
	plugin.DefaultPlugin
}

type KitAuthzConf struct {
	Requirement []Requirement `json:"requirement"`
}

type Requirement struct {
	Namespace string `json:"namespace"`
	Resource  string `json:"resource"`
	Action    string `json:"action"`
}

func (p *KitAuthz) Name() string {
	return "kit_authz"
}

func (p *KitAuthz) ParseConf(in []byte) (interface{}, error) {
	conf := KitAuthzConf{}
	err := json.Unmarshal(in, &conf)
	return conf, err
}

func (p *KitAuthz) RequestFilter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
	requirement := conf.(KitAuthzConf).Requirement
	log.Infof("authz check requirements:%s", strings.Join(lo.Map(requirement, func(t Requirement, _ int) string {
		return fmt.Sprintf("%s/%s@%s", t.Namespace, t.Resource, t.Action)
	}), ","))

	if len(requirement) == 0 {
		return
	}

	ctx := r.Context()
	headers := api.HeaderCarrier(map[string]string{})
	for s, ss := range r.Header().View() {
		headers.Set(s, ss[0])
	}
	for _, contributor := range apiOpt.Propagators {
		//do not handle error
		ctx, _ = contributor.Extract(ctx, headers)
	}
	requirements := lo.Map(requirement, func(t Requirement, _ int) *authz.Requirement {
		return authz.NewRequirement(authz.NewEntityResource(t.Namespace, t.Resource), authz.ActionStr(t.Action))
	})
	resultList, err := authService.BatchCheck(ctx, requirements)
	if err != nil {
		abortWithError(r, err, w)
		return
	}
	for _, result := range resultList {
		if !result.Allowed {
			abortWithError(r, authService.FormatError(ctx, requirements, result), w)
			return
		}
	}
	//authz pass
}
