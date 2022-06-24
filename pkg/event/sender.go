package event

import (
	"context"
	"io"
	"sync"
)

type Sender interface {
	io.Closer
	Send(ctx context.Context, msg Event) error
	BatchSend(ctx context.Context, msg []Event) error
}

type SenderWrap struct {
	Sender
	mu  sync.RWMutex
	mws []SendMiddlewareFunc
}

func (s *SenderWrap) Close() error {
	return s.Close()
}

type SendMiddlewareFunc func(HandlerOf[any]) HandlerOf[any]

func ChainSender(m ...SendMiddlewareFunc) SendMiddlewareFunc {
	return func(next HandlerOf[any]) HandlerOf[any] {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

// NewSender create a *SendWrap with middleware ability
func NewSender(next Sender) *SenderWrap {
	ret := &SenderWrap{
		Sender: next,
	}
	ret.Use(ret.wrapContext())
	return ret
}

func (s *SenderWrap) Use(m ...SendMiddlewareFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, mf := range m {
		s.mws = append(s.mws, mf)
	}
}

func (s *SenderWrap) Send(ctx context.Context, msg Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c := ChainSender(s.mws...)
	h := c(HandlerFuncOf[any](func(ctx context.Context, e any) error {
		// put the real send as inner
		return s.Sender.Send(ctx, e.(Event))
	}))
	return h.Process(ctx, msg)
}

func (s *SenderWrap) BatchSend(ctx context.Context, msg []Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c := ChainSender(s.mws...)
	h := c(HandlerFuncOf[any](func(ctx context.Context, e any) error {
		return s.Sender.BatchSend(ctx, e.([]Event))
	}))
	return h.Process(ctx, msg)
}

func (s *SenderWrap) wrapContext() SendMiddlewareFunc {
	return func(next HandlerOf[any]) HandlerOf[any] {
		return HandlerFuncOf[any](func(ctx context.Context, e any) error {
			ctx = NewSenderContext(ctx, s.Sender)
			return next.Process(ctx, e)
		})
	}
}
