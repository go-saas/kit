package data

import (
	"context"
)

func (a *AggRoot) AfterInsert(ctx context.Context) error {
	return dispatchEvents(ctx, a)
}
func (a *AggRoot) AfterUpdate(ctx context.Context) error {
	return dispatchEvents(ctx, a)
}
func (a *AggRoot) AfterUpsert(ctx context.Context) error {
	return dispatchEvents(ctx, a)
}
