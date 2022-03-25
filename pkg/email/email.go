package email

import (
	"context"
	"crypto/tls"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/lazy"
	"github.com/xhit/go-simple-mail/v2"
	"time"
)

//NewLazyClient create lazy client for email . TODO provider?
func NewLazyClient(cfg *conf.Endpoints) *lazy.Of[*mail.SMTPClient] {
	if cfg.Email == nil || cfg.Email.Smtp == nil {
		panic("endpoints.email.smtp required")
	}
	return lazy.New[*mail.SMTPClient](func(ctx context.Context) (*mail.SMTPClient, error) {

		smtp := cfg.Email.Smtp
		server := mail.NewSMTPClient()

		// SMTP Server
		server.Host = smtp.Host
		server.Port = int(smtp.Port)
		server.Username = smtp.Username
		server.Password = smtp.Password
		switch smtp.Encryption {
		case conf.Email_SMTP_NONE:
			server.Encryption = mail.EncryptionNone
		case conf.Email_SMTP_SSL:
			server.Encryption = mail.EncryptionSSL
		case conf.Email_SMTP_STARTTLS:
			server.Encryption = mail.EncryptionSTARTTLS
		}

		// Variable to keep alive connection
		server.KeepAlive = smtp.KeepAlive
		if smtp.ConnectTimeout != nil {
			// Timeout for connect to SMTP Server
			server.ConnectTimeout = time.Second * time.Duration(smtp.ConnectTimeout.Value)
		}
		if smtp.SendTimeout != nil {
			// Timeout for connect to SMTP Server
			server.SendTimeout = time.Second * time.Duration(smtp.SendTimeout.Value)
		}

		// Set TLSConfig to provide custom TLS configuration. For example,
		// to skip TLS verification (useful for testing):
		server.TLSConfig = &tls.Config{InsecureSkipVerify: smtp.TlsSkipVerify}
		return server.Connect()
	})
}
