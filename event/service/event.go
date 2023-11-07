package service

import (
	"context"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/pkg/api"
	"github.com/samber/lo"

	pb "github.com/go-saas/kit/event/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EventService struct {
	server   *event.ConsumerFactoryServer
	trusted  api.TrustedContextValidator
	producer event.Producer
}

var _ pb.EventServiceServer = (*EventService)(nil)

func NewEventService(server *event.ConsumerFactoryServer, producer event.Producer, trusted api.TrustedContextValidator) *EventService {
	return &EventService{server: server, trusted: trusted, producer: producer}
}

func (s *EventService) HandleEvent(ctx context.Context, req *pb.HandleEventRequest) (*emptypb.Empty, error) {
	if err := api.ErrIfUntrusted(ctx, s.trusted); err != nil {
		return nil, err
	}
	e := req.Message.ToEvent()
	//dispatch to event server
	err := s.server.Process(ctx, e)
	if err != nil {
		return nil, err
	}
	//event handler successfully
	return &emptypb.Empty{}, nil
}

func (s *EventService) PublishEvent(ctx context.Context, req *pb.PublishEventRequest) (*emptypb.Empty, error) {
	if err := api.ErrIfUntrusted(ctx, s.trusted); err != nil {
		return nil, err
	}
	events := lo.Map(req.Messages, func(t *pb.MessageProto, _ int) event.Event {
		return t.ToEvent()
	})
	err := s.producer.BatchSend(ctx, events)
	if err != nil {
		return nil, err
	}
	//produce successfully
	return &emptypb.Empty{}, nil
}
