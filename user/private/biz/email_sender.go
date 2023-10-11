package biz

import (
	"context"
	"fmt"
	kconf "github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/email"
	mail "github.com/wneessen/go-mail"
)

type EmailSender interface {
	//SendForgetPassword send forget password token
	SendForgetPassword(ctx context.Context, email, token string) error
	//SendInviteTenant send invite people into tenant
	SendInviteTenant(ctx context.Context, email, token string) error
	SendPasswordlessLogin(ctx context.Context, email, token string) error
}

// DefaultEmailSender TODO template?
type DefaultEmailSender struct {
	emailer email.Client
	cfg     *kconf.Data
}

func NewEmailSender(emailer email.Client, cfg *kconf.Data) EmailSender {
	return &DefaultEmailSender{emailer: emailer, cfg: cfg}
}

var _ EmailSender = (*DefaultEmailSender)(nil)

func (d *DefaultEmailSender) SendForgetPassword(ctx context.Context, email, token string) error {
	// New email simple html with inline and CC
	// TODO template
	e := mail.NewMsg()
	e.From(d.cfg.Endpoints.Email.From)
	e.To(email)
	e.Subject("Forget Password")
	body := fmt.Sprintf("token: %s", token)
	e.SetBodyString(mail.TypeTextPlain, body)
	err := d.emailer.Send(ctx, e)
	if err != nil {
		return err
	}
	return nil
}

func (d *DefaultEmailSender) SendInviteTenant(ctx context.Context, email, token string) error {
	//TODO implement me
	panic("implement me")
}

func (d *DefaultEmailSender) SendPasswordlessLogin(ctx context.Context, email, token string) error {
	// TODO template
	e := mail.NewMsg()
	e.From(d.cfg.Endpoints.Email.From)
	e.To(email)
	e.Subject("Login")
	body := fmt.Sprintf("token: %s", token)
	e.SetBodyString(mail.TypeTextPlain, body)
	err := d.emailer.Send(ctx, e)
	if err != nil {
		return err
	}
	return nil
}
