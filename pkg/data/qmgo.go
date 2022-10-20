package data

import "context"

type QmgoBeforeInsert interface {
	BeforeInsert(ctx context.Context) error
}
type QmgoAfterInsert interface {
	AfterInsert(ctx context.Context) error
}

type QmgoBeforeUpdate interface {
	BeforeUpdate(ctx context.Context) error
}

type QmgoAfterUpdate interface {
	AfterUpdate(ctx context.Context) error
}

type QmgoBeforeUpsert interface {
	BeforeUpsert(ctx context.Context) error
}

type QmgoAfterUpsert interface {
	AfterUpsert(ctx context.Context) error
}

type QmgoBeforeRemove interface {
	BeforeRemove(ctx context.Context) error
}

type QmgoAfterRemove interface {
	AfterRemove(ctx context.Context) error
}
