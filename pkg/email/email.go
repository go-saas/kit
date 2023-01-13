package email

import (
	"context"
	"fmt"
	"github.com/goava/di"
	mail "github.com/wneessen/go-mail"
	"google.golang.org/protobuf/proto"
)

type Client interface {
	Send(context.Context, ...*mail.Msg) error
}

func NewClient(config *Config, container *di.Container) (Client, error) {
	defConf := &Config{Provider: "log"}
	if config != nil {
		config.Normalize()
		proto.Merge(defConf, config)
	}

	_typeMux.RLock()
	defer _typeMux.RUnlock()

	f, ok := _typeRegister[defConf.Provider]
	if !ok {
		return nil, fmt.Errorf("email provider %s not registered", defConf.Provider)
	}
	return f(config, container)
}
