package biz

import (
	"context"
	"fmt"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	v1 "github.com/goxiaoy/go-saas-kit/saas/event/v1"
)

type TenantReadyEventHandler event.Handler

func NewTenantReadyEventHandler(useCase *TenantUseCase) TenantReadyEventHandler {
	msg := &v1.TenantReadyEvent{}
	return TenantReadyEventHandler(event.ProtoHandler[*v1.TenantReadyEvent](msg, func(ctx context.Context, msg *v1.TenantReadyEvent) error {
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
		return useCase.Update(ctx, tenant, nil)
	}))
}
