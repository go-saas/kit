package service

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	event2 "github.com/go-saas/kit/event"
	v1 "github.com/go-saas/kit/order/api/order/v1"
	pb "github.com/go-saas/kit/payment/api/gateway/v1"
	v13 "github.com/go-saas/kit/payment/event/v1"
	"github.com/go-saas/kit/payment/private/biz"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/authz/authz"
	kithttp "github.com/go-saas/kit/pkg/server/http"
	stripe2 "github.com/go-saas/kit/pkg/stripe"
	"github.com/go-saas/kit/pkg/utils"
	v12 "github.com/go-saas/kit/user/api/user/v1"
	"github.com/stripe/stripe-go/v76"
	stripeclient "github.com/stripe/stripe-go/v76/client"
	"github.com/stripe/stripe-go/v76/webhook"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"strings"
	"time"
)

type PaymentService struct {
	trust            sapi.TrustedContextValidator
	auth             authz.Service
	orderInternalSrv v1.OrderInternalServiceServer
	userInternalSrv  v12.UserInternalServiceServer
	subsRepo         biz.SubscriptionRepo
	stripeClient     *stripeclient.API
	l                *log.Helper
	c                *stripe2.Conf
}

var _ pb.PaymentGatewayServiceServer = (*PaymentService)(nil)
var _ pb.StripePaymentGatewayServiceServer = (*PaymentService)(nil)

func NewPaymentService(
	trust sapi.TrustedContextValidator,
	auth authz.Service,
	orderInternalSrv v1.OrderInternalServiceServer,
	userInternalSrv v12.UserInternalServiceServer,
	subsRepo biz.SubscriptionRepo,
	stripeClient *stripeclient.API,
	logger log.Logger,
	c *stripe2.Conf,
) *PaymentService {
	return &PaymentService{
		trust:            trust,
		auth:             auth,
		orderInternalSrv: orderInternalSrv,
		userInternalSrv:  userInternalSrv,
		subsRepo:         subsRepo,
		stripeClient:     stripeClient,
		l:                log.NewHelper(logger),
		c:                c,
	}
}

func (s *PaymentService) GetPaymentMethod(ctx context.Context, req *pb.GetPaymentMethodRequest) (*pb.GetPaymentMethodReply, error) {
	return &pb.GetPaymentMethodReply{}, nil
}

func (s *PaymentService) GetStripeConfig(ctx context.Context, req *pb.GetStripeConfigRequest) (*pb.GetStripeConfigReply, error) {
	ui, _ := authn.FromUserContext(ctx)
	ret := &pb.GetStripeConfigReply{IsTest: s.c.IsTest, PublishKey: s.c.PublishKey, PriceTables: s.c.PriceTables}
	uid := ui.GetId()
	if len(uid) > 0 {
		customer, err := s.userInternalSrv.FindOrCreateStripeCustomer(ctx, &v12.FindOrCreateStripeCustomerRequest{UserId: &uid})
		if err != nil {
			return nil, err
		}
		ret.CustomerId = customer.StripeCustomerId
	}
	return ret, nil
}

func (s *PaymentService) StripeWebhook(ctx context.Context, req *emptypb.Empty) (*pb.StripeWebhookReply, error) {
	if req, ok := kithttp.ResolveHttpRequest(ctx); ok {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		event, err := webhook.ConstructEvent(body, req.Header.Get("Stripe-Signature"), s.c.WebhookKey)
		if err != nil {
			return nil, handleStripeError(err)
		}
		data := event.Data
		eventType := event.Type
		log.Infof("receive event type %s with data %v", eventType, data.Raw)
		switch eventType {
		case "customer.subscription.updated":
			subs := &stripe.Subscription{}
			err = json.Unmarshal(data.Raw, subs)
			if err != nil {
				return nil, errors.BadRequest("", "")
			}
			localSubs, err := s.subsRepo.FindByProvider(ctx, stripe2.ProviderName, subs.ID)
			if err != nil {
				return nil, err
			}
			MapStripeSubscription2Biz(subs, localSubs)
			ee, _ := event2.NewMessageFromProto(&v13.SubscriptionChangedEvent{Id: localSubs.ID.String()})
			localSubs.AppendEvent(ee)
			err = s.subsRepo.Update(ctx, localSubs.ID.String(), localSubs, nil)
			if err != nil {
				return nil, err
			}
		case "invoice.payment_succeeded":
			var invoice stripe.Invoice
			err = json.Unmarshal(event.Data.Raw, &invoice)
			if err != nil {
				return nil, errors.BadRequest("", "")
			}
			pi, _ := s.stripeClient.PaymentIntents.Get(
				invoice.PaymentIntent.ID,
				nil,
			)
			//set default payment method
			params := &stripe.SubscriptionParams{
				DefaultPaymentMethod: stripe.String(pi.PaymentMethod.ID),
			}
			_, err = s.stripeClient.Subscriptions.Update(invoice.Subscription.ID, params)
			if err != nil {
				return nil, err
			}
		case "payment_intent.succeeded":
			intent := stripe.PaymentIntent{}
			err = json.Unmarshal(data.Raw, &intent)
			if err != nil {
				return nil, errors.BadRequest("", "")
			}
			orderId := intent.Metadata["order_id"]
			//no order id, maybe subscription,handled by invoice
			// has order id, one time purchase
			if len(orderId) != 0 {
				t := time.Now()
				_, err = s.orderInternalSrv.InternalOrderPaySuccess(ctx, &v1.InternalOrderPaySuccessRequest{
					Id:              orderId,
					PayExtra:        utils.Map2Structpb(data.Object),
					PaidPriceAmount: intent.Amount,
					CurrencyCode:    strings.ToUpper(string(intent.Currency)),
					PayProvider:     stripe2.ProviderName,
					PaidTime:        utils.Time2Timepb(&t),
				})
				if err != nil {
					return nil, err
				}
			}
		case "charge.refunded":
			refund := stripe.Refund{}
			json.Unmarshal(data.Raw, &refund)
			intent, err := s.stripeClient.PaymentIntents.Get(refund.PaymentIntent.ID, nil)
			if err != nil {
				return nil, handleStripeError(err)
			}
			t := time.Now()
			_, err = s.orderInternalSrv.InternalOrderRefunded(ctx, &v1.InternalOrderRefundedRequest{
				Id:                intent.Metadata["order_id"],
				PayExtra:          utils.Map2Structpb(data.Object),
				RefundTime:        utils.Time2Timepb(&t),
				RefundPriceAmount: refund.Amount,
				CurrencyCode:      strings.ToUpper(string(refund.Currency)),
				PayProvider:       stripe2.ProviderName,
			})
			if err != nil {
				return nil, err
			}
		case "payment_intent.payment_failed":
		case "setup_intent.setup_failed":
		case "setup_intent.succeeded":
		case "setup_intent.created":
		default:

		}
		return &pb.StripeWebhookReply{}, nil
	} else {
		return nil, errors.BadRequest("", "")
	}
}

func handleStripeError(err error) error {
	//TODO handle stripe
	return err
}
