package uow

import "context"

type UnitOfWorkKey string

type CancelFunc func(ctx context.Context) context.Context
