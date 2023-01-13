package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-saas/kit/pkg/email"
	"github.com/goava/di"
	mail "github.com/wneessen/go-mail"
)

func init() {
	email.RegisterProvider("smtp", func(config *email.Config, container *di.Container) (email.Client, error) {
		if config.Smtp == nil {
			return nil, fmt.Errorf("smtp config is required")
		}
		smtp := config.Smtp
		var opts []mail.Option

		if smtp.Port != nil {
			opts = append(opts, mail.WithPort(int(*smtp.Port)))
		}
		opts = append(opts, mail.WithUsername(smtp.Username))
		opts = append(opts, mail.WithPassword(smtp.Password))

		opts = append(opts, mail.WithTLSConfig(&tls.Config{InsecureSkipVerify: smtp.TlsSkipVerify}))

		if smtp.Timeout != nil {
			opts = append(opts, mail.WithTimeout(smtp.Timeout.AsDuration()))
		}
		c, err := mail.NewClient(smtp.Host, opts...)
		if err != nil {
			return nil, err
		}
		return &client{up: c}, nil

	})
}

type client struct {
	up *mail.Client
}

var _ email.Client = (*client)(nil)

func (c *client) Send(ctx context.Context, email ...*mail.Msg) error {
	return c.up.DialAndSendWithContext(ctx, email...)

}
