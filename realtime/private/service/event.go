package service

import (
	"context"
	"fmt"
	"github.com/centrifugal/centrifuge"
	"github.com/go-saas/kit/event"
	v1 "github.com/go-saas/kit/realtime/event/v1"
	"github.com/go-saas/kit/realtime/private/biz"
	"google.golang.org/grpc/encoding"
)

func NewNotificationEventHandler(node *centrifuge.Node, repo biz.NotificationRepo) event.ConsumerHandler {
	msg := &v1.NotificationEvent{}
	return event.ProtoHandler[*v1.NotificationEvent](msg, event.HandlerFuncOf[*v1.NotificationEvent](
		func(ctx context.Context, msg *v1.NotificationEvent) error {
			//store
			notifications := biz.FromNotificationEvents(msg)
			err := repo.BatchCreate(ctx, notifications, 100)
			if err != nil {
				return err
			}
			//publish to channel
			for _, notification := range notifications {
				r := MapBizNotification2Pb(notification, 0)
				data, _ := encoding.GetCodec("json").Marshal(r)
				_, err := node.Publish(fmt.Sprintf("notification#%s", notification.UserId), data)
				if err != nil {
					return err
				}
			}
			return err
		}))
}
