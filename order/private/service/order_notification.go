package service

import (
	"context"
	"github.com/go-saas/kit/event"
	v13 "github.com/go-saas/kit/order/api/order/v1"
	v1 "github.com/go-saas/kit/order/event/v1"
	v12 "github.com/go-saas/kit/realtime/event/v1"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"strings"
)

func NewOrderSuccessNotification(producer event.Producer) event.ConsumerHandler {
	msg := &v1.OrderPaySuccessEvent{}
	return event.ProtoHandler[*v1.OrderPaySuccessEvent](msg, event.HandlerFuncOf[*v1.OrderPaySuccessEvent](func(ctx context.Context, msg *v1.OrderPaySuccessEvent) error {
		itemNames := lo.Map(msg.Order.Items, func(item *v13.OrderItem, _ int) string {
			return strings.Join([]string{item.Product.Name, item.Product.SkuTitle}, " ")
		})
		oData, _ := protojson.Marshal(msg.Order)
		extra := new(structpb.Struct)
		protojson.Unmarshal(oData, extra)

		notification := &v12.NotificationEvent{
			Title:   "Order Paid!",
			Level:   v12.NotificationLevel_INFO,
			Desc:    strings.Join(itemNames, ","),
			UserIds: []string{msg.Order.CustomerId},
			Extra:   extra,
		}
		ee, _ := event.NewMessageFromProto(notification)
		return producer.Send(ctx, ee)
	}))
}
