package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/go-saas/kit/order/api/order/v1"
	pb "github.com/go-saas/kit/payment/api/gateway/v1"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/price"
	kithttp "github.com/go-saas/kit/pkg/server/http"
	stripe2 "github.com/go-saas/kit/pkg/stripe"
	"github.com/go-saas/kit/pkg/utils"
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
	stripeClient     *stripeclient.API
	l                *log.Helper
	c                *stripe2.StripeConf
}

var _ pb.PaymentGatewayServiceServer = (*PaymentService)(nil)
var _ pb.StripePaymentGatewayServiceServer = (*PaymentService)(nil)

func NewPaymentService(
	trust sapi.TrustedContextValidator,
	auth authz.Service,
	orderInternalSrv v1.OrderInternalServiceServer,
	stripeClient *stripeclient.API,
	logger log.Logger,
	c *stripe2.StripeConf,
) *PaymentService {
	return &PaymentService{
		trust:            trust,
		auth:             auth,
		orderInternalSrv: orderInternalSrv,
		stripeClient:     stripeClient,
		l:                log.NewHelper(logger),
		c:                c,
	}
}

func (s *PaymentService) GetPaymentMethod(ctx context.Context, req *pb.GetPaymentMethodRequest) (*pb.GetPaymentMethodReply, error) {
	return &pb.GetPaymentMethodReply{}, nil
}

func (s *PaymentService) GetStripeConfig(ctx context.Context, req *pb.GetStripeConfigRequest) (*pb.GetStripeConfigReply, error) {
	return &pb.GetStripeConfigReply{IsTest: s.c.IsTest, PublishKey: s.c.PublishKey}, nil
}

func (s *PaymentService) CreateStripePaymentIntent(ctx context.Context, req *pb.CreateStripePaymentIntentRequest) (*pb.CreateStripePaymentIntentReply, error) {
	userInfo, err := authn.ErrIfUnauthenticated(ctx)
	if err != nil {
		return nil, err
	}
	order, err := s.orderInternalSrv.GetInternalOrder(ctx, &v1.GetInternalOrderRequest{Id: req.OrderId})
	if err != nil {
		return nil, err
	}
	if order.CustomerId != userInfo.GetId() {
		return nil, errors.NotFound("", "")
	}
	userId := userInfo.GetId()
	customer, err := s.findOrCreateCustomer(userId)
	ephemeralKey, err := s.stripeClient.EphemeralKeys.New(&stripe.EphemeralKeyParams{
		Customer:      &customer.ID,
		StripeVersion: stripe.String(stripe.APIVersion),
	})
	if err != nil {
		return nil, handleStripeError(err)
	}

	paymentIntentParams := &stripe.PaymentIntentParams{
		Amount:   &order.TotalPrice.Amount,
		Currency: &order.TotalPrice.CurrencyCode,
		Customer: &customer.ID,
	}
	paymentIntentParams.Metadata = map[string]string{
		"user_id":  userInfo.GetId(),
		"order_id": req.OrderId,
	}
	intent, err := s.stripeClient.PaymentIntents.New(paymentIntentParams)
	if err != nil {
		return nil, handleStripeError(err)
	}
	return &pb.CreateStripePaymentIntentReply{
		PaymentIntent: intent.ClientSecret,
		CustomerId:    customer.ID,
		EphemeralKey:  ephemeralKey.Secret,
	}, nil
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
		case "payment_intent.succeeded":
			intent := stripe.PaymentIntent{}
			json.Unmarshal(data.Raw, &intent)
			totalPrice, err := price.NewPriceFromInt64(intent.Amount, strings.ToUpper(string(intent.Currency)))
			if err != nil {
				return nil, err
			}
			//TODO paid time
			t := time.Now()
			_, err = s.orderInternalSrv.InternalOrderPaySuccess(ctx, &v1.InternalOrderPaySuccessRequest{
				Id:        intent.Metadata["order_id"],
				PayExtra:  utils.Map2Structpb(data.Object),
				PaidPrice: totalPrice.ToPricePb(ctx),
				PayWay:    "stripe",
				PaidTime:  utils.Time2Timepb(&t),
			})
			if err != nil {
				return nil, err
			}
		case "charge.refunded":
			refund := stripe.Refund{}
			json.Unmarshal(data.Raw, &refund)
			refundPrice, err := price.NewPriceFromInt64(refund.Amount, strings.ToUpper(string(refund.Currency)))
			if err != nil {
				return nil, err
			}
			intent, err := s.stripeClient.PaymentIntents.Get(refund.PaymentIntent.ID, nil)
			if err != nil {
				return nil, handleStripeError(err)
			}
			t := time.Now()
			_, err = s.orderInternalSrv.InternalOrderRefunded(ctx, &v1.InternalOrderRefundedRequest{
				Id:          intent.Metadata["order_id"],
				PayExtra:    utils.Map2Structpb(data.Object),
				RefundTime:  utils.Time2Timepb(&t),
				RefundPrice: refundPrice.ToPricePb(ctx),
				PayWay:      "stripe",
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

func (s *PaymentService) findOrCreateCustomer(userId string) (*stripe.Customer, error) {
	var err error
	customerSearch := &stripe.CustomerSearchParams{}
	customerSearch.Query = fmt.Sprintf("metadata['user_id']:'%s'", userId)
	searchIter := s.stripeClient.Customers.Search(customerSearch)
	if searchIter.Err() != nil {
		return nil, handleStripeError(searchIter.Err())
	}
	var customer *stripe.Customer
	for searchIter.Next() {
		if searchIter.Err() != nil {
			return nil, handleStripeError(searchIter.Err())
		}
		customer = searchIter.Customer()
		break
	}
	if searchIter.Err() != nil {
		return nil, handleStripeError(searchIter.Err())
	}
	if customer == nil {
		params := &stripe.CustomerParams{
			Name: &userId,
		}
		params.Metadata = map[string]string{
			"user_id": userId,
		}
		customer, err = s.stripeClient.Customers.New(params)
		if err != nil {
			return nil, handleStripeError(err)
		}
	}
	return customer, nil
}
