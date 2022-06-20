package job

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/hibiken/asynq"
)

type Server struct {
	*asynq.ServeMux
	server LazyAsynqServer
}

var _ transport.Server = (*Server)(nil)

func NewServer(s LazyAsynqServer) *Server {
	return &Server{server: s, ServeMux: asynq.NewServeMux()}
}

func (s *Server) Start(ctx context.Context) error {
	asynqServer, err := s.server.Value(ctx)
	if err != nil {
		return err
	}
	return asynqServer.Start(s)
}

func (s *Server) Stop(ctx context.Context) error {
	asynqServer, err := s.server.Value(ctx)
	if err != nil {
		return err
	}
	asynqServer.Shutdown()
	return nil
}
