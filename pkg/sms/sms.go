package sms

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

// SMSSender sends SMS messages to a phone number
type SMSSender interface {
	Send(ctx context.Context, number, text string) error
}

type SMSLogSender struct {
	log *log.Helper
}

func NewSMSLogSender(l log.Logger) *SMSLogSender {
	return &SMSLogSender{log: log.NewHelper(l)}
}

// Send an SMS
func (s *SMSLogSender) Send(_ context.Context, number, text string) error {
	s.log.Info("sms sent to:", number, "contents:", text)
	return nil
}
