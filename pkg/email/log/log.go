package log

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/pkg/email"
	"github.com/goava/di"
	"github.com/samber/lo"
	mail "github.com/wneessen/go-mail"
)

func init() {
	email.RegisterProvider("log", func(c *email.Config, container *di.Container) (email.Client, error) {
		var logger log.Logger
		err := container.Resolve(&logger)
		if err != nil {
			return nil, err
		}
		return &client{l: log.NewHelper(logger)}, nil
	})
}

type client struct {
	l *log.Helper
}

func (c *client) Send(ctx context.Context, email ...*mail.Msg) error {
	for _, m := range email {
		parts := lo.Map(m.GetParts(), func(t *mail.Part, _ int) string {
			content, _ := t.GetContent()
			return string(content)
		})
		c.l.Infow(log.DefaultMessageKey, "send email", "to", m.GetToString(), "parts", parts)
	}
	return nil
}

var _ email.Client = (*client)(nil)
