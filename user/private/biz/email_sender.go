package biz

import (
	"context"
	"fmt"
	kconf "github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/email"
	mail "github.com/xhit/go-simple-mail/v2"
)

type EmailSender interface {
	//SendForgetPassword send forget password token
	SendForgetPassword(ctx context.Context, email, token string) error
	//SendInviteTenant send invite people into tenant
	SendInviteTenant(ctx context.Context, email, token string) error
}

//DefaultEmailSender TODO template?
type DefaultEmailSender struct {
	emailer email.LazyClient
	cfg     *kconf.Data
}

func NewEmailSender(emailer email.LazyClient, cfg *kconf.Data) EmailSender {
	return &DefaultEmailSender{emailer: emailer, cfg: cfg}
}

var _ EmailSender = (*DefaultEmailSender)(nil)

func (d *DefaultEmailSender) SendForgetPassword(ctx context.Context, email, token string) error {
	// New email simple html with inline and CC
	e := mail.NewMSG()
	e.SetFrom(d.cfg.Endpoints.Email.From).
		AddTo(email).
		SetSubject("Forget Password")
	body := fmt.Sprintf("token: %s", token)
	e.SetBody(mail.TextPlain, body)
	if e.Error != nil {
		return e.Error
	}
	client, err := d.emailer.Value(ctx)
	if err != nil {
		return err
	}
	err = e.Send(client)
	if err != nil {
		return err
	}
	return nil
}

func (d *DefaultEmailSender) SendInviteTenant(ctx context.Context, email, token string) error {
	//TODO implement me
	panic("implement me")
}
