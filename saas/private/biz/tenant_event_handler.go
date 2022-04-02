package biz

import (
	"context"
	"fmt"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	v1 "github.com/goxiaoy/go-saas-kit/saas/event/v1"
	"github.com/samber/lo"
)

type TenantReadyEventHandler event.Handler

func NewTenantReadyEventHandler(useCase *TenantUseCase) TenantReadyEventHandler {
	msg := &v1.TenantReadyEvent{}
	return TenantReadyEventHandler(event.ProtoHandler[*v1.TenantReadyEvent](msg, func(ctx context.Context, msg *v1.TenantReadyEvent) error {
		tenant, err := useCase.FindByIdOrName(ctx, msg.Id)
		if err != nil {
			return err
		}
		if len(msg.ConnStrKey) > 0 {
			tenant.Conn = lo.Filter(tenant.Conn, func(conn TenantConn, _ int) bool { return conn.Key != msg.ConnStrKey })
			tenant.Conn = append(tenant.Conn, TenantConn{
				Key:   msg.ConnStrKey,
				Value: msg.ConnStrValue,
			})
		}
		if tenant.Extra == nil {
			tenant.Extra = map[string]interface{}{}
		}
		if len(msg.ServiceName) > 0 {
			tenant.Extra[fmt.Sprintf("%s_status", msg.ServiceName)] = "READY"
		}
		return useCase.Update(ctx, tenant, nil)
	}))
}
