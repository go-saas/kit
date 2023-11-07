package biz

import (
	"context"
	"fmt"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/event"
	v15 "github.com/go-saas/kit/order/api/order/v1"
	v13 "github.com/go-saas/kit/order/event/v1"
	v14 "github.com/go-saas/kit/payment/api/subscription/v1"
	v12 "github.com/go-saas/kit/payment/event/v1"
	"github.com/go-saas/kit/payment/private/biz"
	"github.com/go-saas/kit/pkg/query"
	v1 "github.com/go-saas/kit/saas/event/v1"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func NewTenantReadyEventHandler(useCase *TenantUseCase) event.ConsumerHandler {
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

func NewSubscriptionChangedEventHandler(useCase *TenantUseCase, planRepo PlanRepo, subsSrv v14.SubscriptionInternalServiceServer) event.ConsumerHandler {
	msg := &v12.SubscriptionChangedEvent{}
	return event.ProtoHandler[*v12.SubscriptionChangedEvent](msg, event.HandlerFuncOf[*v12.SubscriptionChangedEvent](func(ctx context.Context, msg *v12.SubscriptionChangedEvent) error {
		klog.Infof("receive msg SubscriptionChangedEvent")
		subs, err := subsSrv.GetInternalSubscription(ctx, &v14.GetInternalSubscriptionRequest{Id: msg.GetId()})
		if err != nil {
			return err
		}
		var plan *Plan
		var tenantId string
		for _, item := range subs.Items {
			if item.BizPayload != nil {
				t, ok := item.BizPayload.AsMap()["tenant_id"].(string)
				if ok {
					tenantId = t
				}
			}
			plan, err = planRepo.FindByProductId(ctx, item.ProductId)
			if err != nil {
				return err
			}
			if plan != nil {
				break
			}
		}
		if plan != nil && len(tenantId) > 0 {
			tenantEntity, err := useCase.Get(ctx, tenantId)
			if err != nil {
				return err
			}
			if subs.Status == string(biz.SubscriptionStatusActive) {
				tenantEntity.PlanKey = &plan.Key
				tenantEntity.ActiveSubscriptionID = &subs.Id
				err = useCase.Update(ctx, tenantEntity, query.NewField(&fieldmaskpb.FieldMask{Paths: []string{"plan_key", "active_subscription_id"}}))
				if err != nil {
					return err
				}
			}
		}
		return nil
	}))
}

func NewOrderChangedEventHandler(useCase *TenantUseCase, planRepo PlanRepo, orderSrv v15.OrderInternalServiceServer) event.ConsumerHandler {
	msg := &v13.OrderPaySuccessEvent{}
	return event.ProtoHandler[*v13.OrderPaySuccessEvent](msg, event.HandlerFuncOf[*v13.OrderPaySuccessEvent](func(ctx context.Context, msg *v13.OrderPaySuccessEvent) error {
		//TODO change plan
		klog.Infof("receive msg OrderPaySuccessEvent")
		order, err := orderSrv.GetInternalOrder(ctx, &v15.GetInternalOrderRequest{Id: &msg.Order.Id})
		if err != nil {
			return err
		}
		var plan *Plan
		var tenantId string
		for _, item := range order.Items {
			if item.BizPayload != nil {
				t, ok := item.BizPayload.AsMap()["tenant_id"].(string)
				if ok {
					tenantId = t
				}
			}
			if item.Product == nil || item.Product.Id == nil {
				continue
			}
			plan, err = planRepo.FindByProductId(ctx, *item.Product.Id)
			if err != nil {
				return err
			}
			if plan != nil {
				break
			}
		}
		if plan != nil && len(tenantId) > 0 {
			tenantEntity, err := useCase.Get(ctx, tenantId)
			if err != nil {
				return err
			}
			tenantEntity.PlanKey = &plan.Key
			err = useCase.Update(ctx, tenantEntity, query.NewField(&fieldmaskpb.FieldMask{Paths: []string{"plan_key", "active_subscription_id"}}))
			if err != nil {
				return err
			}
		}
		return nil
	}))
}
