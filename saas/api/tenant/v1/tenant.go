package v1

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"regexp"
)

func (x *TenantInfo) NormalizeHost(ctx context.Context, app *conf.AppConfig) {
	x.Host = normalizeHost(ctx, app, x.Name)
}

func (x *Tenant) NormalizeHost(ctx context.Context, app *conf.AppConfig) {
	x.Host = normalizeHost(ctx, app, x.Name)
}

func normalizeHost(ctx context.Context, app *conf.AppConfig, name string) string {
	if app == nil || app.DomainFormat == nil {
		return ""
	}
	if t, ok := transport.FromServerContext(ctx); ok {
		if ht, ok := t.(*khttp.Transport); ok {
			host := ht.Request().Host
			if len(host) == 0 {
				return ""
			}
			reg := regexp.MustCompile(app.DomainFormat.Value)
			return reg.ReplaceAllString(host, name)
		}
	}
	return ""
}
