package event

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
)

// FactoryServer resolve LazyReceiver from factory, then wrap as kratos server
type FactoryServer struct {
	*ServeMux
	lr  LazyReceiver
	cfg *conf.Event
}

var _ transport.Server = (*FactoryServer)(nil)

func (f *FactoryServer) Start(ctx context.Context) error {
	if f.r != nil {
		panic("server can not start twice")
	}
	r, err := f.lr(ctx, f.cfg)
	if err != nil {
		return err
	}
	f.r = r
	return r.Receive(ctx, f)
}

func (f *FactoryServer) Stop(ctx context.Context) error {
	return f.r.Close()
}

type Server struct {
	*ServeMux
}

var _ transport.Server = (*Server)(nil)

// NewServer create server from Receiver directly
func NewServer(r Receiver) *Server {
	return &Server{&ServeMux{r: r}}
}

func (s *Server) Start(ctx context.Context) error {
	return s.r.Receive(ctx, s)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.r.Close()
}
