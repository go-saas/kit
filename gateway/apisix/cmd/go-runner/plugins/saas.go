package plugins

import (
	"context"
	"encoding/base64"
	"encoding/json"
	pkgHTTP "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/goxiaoy/go-saas-kit/pkg/api"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	"net/http"
)

func Saas(ctx context.Context, tenantStore common.TenantStore, pathRegex string, w http.ResponseWriter, r pkgHTTP.Request) (context.Context, error) {

	key := shttp.KeyOrDefault("")

	//get tenant config
	tenantConfigProvider := common.NewDefaultTenantConfigProvider(NewResolver(r, key, pathRegex), tenantStore)
	tenantConfig, ctx, err := tenantConfigProvider.Get(ctx)
	if err != nil {
		return ctx, err
	}

	//extract previous id or name for logging
	resolveValue := common.FromTenantResolveRes(ctx)
	idOrName := ""
	if resolveValue != nil {
		idOrName = resolveValue.TenantIdOrName
	}
	log.Infof("resolve raw tenant: %s , id: %s ,is host: %v", idOrName, tenantConfig.ID, len(tenantConfig.ID) == 0)
	r.Header().Set(api.TenantKey, tenantConfig.ID)

	b, _ := json.Marshal(tenantConfig)
	r.Header().Set(api.TenantInfoKey, base64.StdEncoding.EncodeToString(b))
	return common.NewCurrentTenant(ctx, tenantConfig.ID, tenantConfig.Name), nil
}
