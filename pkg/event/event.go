package event

import (
	"context"
	"errors"
	"fmt"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	"github.com/goxiaoy/go-saas-kit/pkg/event/kafka"
	"strings"
)

func NewEventReceiver(cfg *conf.Event, handler event.Handler) (event.Receiver, func(), error) {
	var ret event.Receiver
	var err error
	if cfg.Type == "kafka" || cfg.Type == "" {
		var addr []string
		if cfg.Addr!=""{
			addr = strings.Split(cfg.Addr, ";")
		}else {
			addr = []string{"localhost:9092"}
		}
		ret, err = kafka.NewKafkaReceiver(addr, cfg.Topic, cfg.Group)
	} else {
		return nil, nil, errors.New(fmt.Sprintf("unsupported event type %s", cfg.Type))
	}
	if err != nil {
		return nil, func() {}, err
	}
	err = ret.Receive(context.Background(), handler)
	return ret, func() {
		ret.Close()
	}, err
}
