package plugins

import (
	"encoding/json"
	"fmt"
	pkgHTTP "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas-kit/pkg/authz/authz"
	"github.com/samber/lo"
	"net/http"
	"strings"
)

// KitAuthz is a demo to show how to return data directly instead of proxying
// it to the upstream.
type KitAuthz struct {
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

func (p *KitAuthz) Filter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
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

	subjects, _ := subjectResolver.ResolveFromContext(ctx)
	log.Infof("subjects: %s", strings.Join(lo.Map(subjects, func(t authz.Subject, _ int) string {
		return t.GetIdentity()
	}), ","))

	resultList, err := authService.BatchCheck(ctx, lo.Map(requirement, func(t Requirement, _ int) *authz.Requirement {
		return authz.NewRequirement(authz.NewEntityResource(t.Namespace, t.Resource), authz.ActionStr(t.Action))
	}))
	if err != nil {
		abortWithError(err, w)
		return
	}
	for _, result := range resultList {
		if !result.Allowed {
			abortWithError(authz.FormatError(ctx, result, subjects...), w)
			return
		}
	}
	//authz pass
}
