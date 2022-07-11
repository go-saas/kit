package event

import (
	"context"
	"database/sql"
	"github.com/go-saas/kit/event"
	"github.com/go-saas/uow"
	"sync"
)

const (
	UowKind = "event"
)

type Transactional struct {
	ctx      context.Context
	producer event.Producer
	events   []event.Event
	sync.Mutex
}

func NewTransactional(ctx context.Context, producer event.Producer) *Transactional {
	return &Transactional{
		ctx:      ctx,
		producer: producer,
	}
}

var (
	_ uow.TransactionalDb = (*Transactional)(nil)
	_ uow.Txn             = (*Transactional)(nil)
)

func (t *Transactional) Commit() error {
	if len(t.events) == 0 {
		return nil
	}
	return t.producer.BatchSend(t.ctx, t.events)
}

func (t *Transactional) Rollback() error {
	//can not perform rollback
	return nil
}

func (t *Transactional) Begin(opt ...*sql.TxOptions) (db uow.Txn, err error) {
	return NewTransactional(t.ctx, t.producer), nil
}

func (t *Transactional) Send(msg ...event.Event) error {
	t.Lock()
	defer t.Unlock()
	t.events = append(t.events, msg...)
	return nil
}
