package data

import (
	"context"
	"errors"
	"github.com/go-saas/kit/event"
	event2 "github.com/go-saas/kit/pkg/uow/event"
	"github.com/go-saas/uow"
)

var (
	ErrProducerNotFound = errors.New("producer not found")
)

type Agg interface {
	AppendEvent(events ...event.Event)
	ConsumeEventsIfAny(ctx context.Context, fn func(ctx context.Context, events []event.Event) error) (err error)
}

// AggRoot DDD aggregate root. will dispatch events after transaction committed
type AggRoot struct {
	events []event.Event
}

func (a *AggRoot) AppendEvent(events ...event.Event) {
	a.events = append(a.events, events...)
}

func (a *AggRoot) ConsumeEventsIfAny(ctx context.Context, fn func(ctx context.Context, events []event.Event) error) (err error) {
	//TODO lock?
	if len(a.events) > 0 {
		err = fn(ctx, a.events)
		if err == nil {
			//clear events
			a.events = nil
		}
		return err
	}
	return nil
}

func dispatchEvents(ctx context.Context, agg Agg) error {
	return agg.ConsumeEventsIfAny(ctx, func(ctx context.Context, events []event.Event) error {
		if uow, ok := uow.FromCurrentUow(ctx); ok {
			// uow manage events
			tdb, err := uow.GetTxDb(ctx, event2.UowKind)
			if err != nil {
				return err
			}
			return tdb.(*event2.Transactional).Send(events...)
		} else {
			//send immediately
			if p, ok := event.FromProducerContext(ctx); ok {
				return p.BatchSend(ctx, events)
			} else {
				return ErrProducerNotFound
			}
		}
		return nil
	})
}
