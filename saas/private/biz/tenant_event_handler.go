package biz

import (
	"context"
	"fmt"
	"github.com/go-saas/kit/event"
	v1 "github.com/go-saas/kit/saas/event/v1"
	"github.com/samber/lo"
)

type TenantReadyEventHandler event.ConsumerHandler

func NewTenantReadyEventHandler(useCase *TenantUseCase) TenantReadyEventHandler {
	msg := &v1.TenantReadyEvent{}
	return event.ProtoHandler[*v1.TenantReadyEvent](msg, event.HandlerFuncOf[*v1.TenantReadyEvent](func(ctx context.Context, msg *v1.TenantReadyEvent) error {
		tenant, err := useCase.FindByIdOrName(ctx, msg.Id)
		if err != nil {
			return err
		}

		if tenant.Extra == nil {
			tenant.Extra = map[string]interface{}{}
		}
		if len(msg.ServiceName) > 0 {
			tenant.Extra[fmt.Sprintf("%s_status", msg.ServiceName)] = "READY"
		}
		if c, ok := lo.Find(tenant.Conn, func(c TenantConn) bool { return c.Key == msg.ServiceName }); ok {
			c.Ready = true
		}
		return useCase.Update(ctx, tenant, nil)
	}))
}
