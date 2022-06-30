package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/kit/pkg/api"

	pb "github.com/go-saas/kit/event/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EventService struct {
	pb.UnimplementedEventServiceServer
	server  *event.ConsumerFactoryServer
	trusted api.TrustedContextValidator
}

func NewEventService(server *event.ConsumerFactoryServer, trusted api.TrustedContextValidator) *EventService {
	return &EventService{server: server, trusted: trusted}
}

func (s *EventService) Event(ctx context.Context, req *pb.EventRequest) (*emptypb.Empty, error) {
	if ok, err := s.trusted.Trusted(ctx); err != nil {
		return nil, err
	} else if ok {
		e := req.Message.ToEvent()
		//dispatch to event server
		err := s.server.Process(ctx, e)
		if err != nil {
			return nil, err
		}
		//event handler successfully
		return &emptypb.Empty{}, nil
	} else {
		// internal api only
		return nil, errors.Forbidden("", "")
	}
}
