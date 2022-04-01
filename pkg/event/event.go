package event

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	"github.com/goxiaoy/go-saas-kit/pkg/event/kafka"
	"strings"
)

func NewEventReceiver(cfg *conf.Event, handler event.Handler, logger log.Logger) (event.Receiver, func(), error) {
	var ret event.Receiver
	var err error
	var clean = func() {}
	if cfg.Type == "kafka" || cfg.Type == "" {
		var addr []string
		if cfg.Addr != "" {
			addr = strings.Split(cfg.Addr, ";")
		} else {
			addr = []string{"localhost:9092"}
		}
		ret, clean, err = kafka.NewKafkaReceiver(addr, cfg.Topic, cfg.Group, logger, cfg.Kafka)
	} else {
		return nil, nil, errors.New(fmt.Sprintf("unsupported event type %s", cfg.Type))
	}
	if err != nil {
		return nil, func() {}, err
	}
	err = ret.Receive(context.Background(), handler)
	return ret, clean, err
}

func NewEventSender(cfg *conf.Event, logger log.Logger) (event.Sender, func(), error) {
	var ret event.Sender
	var err error
	var clean = func() {}
	if cfg.Type == "kafka" || cfg.Type == "" {
		var addr []string
		if cfg.Addr != "" {
			addr = strings.Split(cfg.Addr, ";")
		} else {
			addr = []string{"localhost:9092"}
		}
		ret, clean, err = kafka.NewKafkaSender(addr, cfg.Topic, logger, cfg.Kafka)
	} else {
		return nil, nil, errors.New(fmt.Sprintf("unsupported event type %s", cfg.Type))
	}
	if err != nil {
		return nil, clean, err
	}
	return ret, clean, err
}
