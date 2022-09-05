package server

import (
	"context"
	"github.com/centrifugal/centrifuge"
	"github.com/go-kratos/kratos/v2/transport"
)

type Centrifuge struct {
	node *centrifuge.Node
}

func NewCentrifuge(node *centrifuge.Node) *Centrifuge {
	return &Centrifuge{node: node}
}

var _ transport.Server = (*Centrifuge)(nil)

func (c *Centrifuge) Start(ctx context.Context) error {
	return c.node.Run()
}

func (c *Centrifuge) Stop(ctx context.Context) error {
	return c.node.Shutdown(ctx)
}
