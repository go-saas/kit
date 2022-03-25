package sms

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

// Sender sends SMS messages to a phone number
type Sender interface {
	Send(ctx context.Context, number, text string) error
}

type LogSender struct {
	log *log.Helper
}

func NewSMSLogSender(l log.Logger) *LogSender {
	return &LogSender{log: log.NewHelper(l)}
}

// Send an SMS
func (s *LogSender) Send(_ context.Context, number, text string) error {
	s.log.Info("sms sent to:", number, "contents:", text)
	return nil
}
