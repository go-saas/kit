package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/event"
	"github.com/goxiaoy/go-saas-kit/pkg/event/trace"
	"strings"
	"sync"
)

var (
	_ event.Sender   = (*kafkaSender)(nil)
	_ event.Receiver = (*kafkaReceiver)(nil)
)

func init() {
	event.RegisterReceiver("kafka", func(ctx context.Context, cfg *conf.Event) (event.Receiver, error) {
		var addr []string
		if cfg.Addr != "" {
			addr = strings.Split(cfg.Addr, ";")
		} else {
			addr = []string{"localhost:9092"}
		}
		return NewKafkaReceiver(addr, cfg.Topic, cfg.Group, cfg.Kafka)
	})

	event.RegisterSender("kafka", func(cfg *conf.Event) (event.Sender, func(), error) {
		var addr []string
		if cfg.Addr != "" {
			addr = strings.Split(cfg.Addr, ";")
		} else {
			addr = []string{"localhost:9092"}
		}
		sender, c, err := NewKafkaSender(addr, cfg.Topic, cfg.Kafka)
		if err != nil {
			return nil, c, err
		}
		res := event.NewSender(sender)
		res.Use(trace.Send())
		return res, c, nil
	})
}

type kafkaSender struct {
	writer  sarama.SyncProducer
	topic   string
	address []string
}

func NewKafkaSender(address []string, topic string, cfg *conf.Event_Kafka) (*kafkaSender, func(), error) {
	conf := sarama.NewConfig()

	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true
	var cleanup = func() {

	}
	err := patchConf(conf, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	w, err := sarama.NewSyncProducer(address, conf)
	if err != nil {
		return nil, cleanup, err
	}
	s := &kafkaSender{writer: w, topic: topic, address: address}
	return s, func() {
		s.Close()
	}, nil
}

func (s *kafkaSender) Send(ctx context.Context, message event.Event) error {
	_, _, err := s.writer.SendMessage(s.toMsg(ctx, message))
	if err != nil {
		return err
	}
	return nil
}

func (s *kafkaSender) BatchSend(ctx context.Context, message []event.Event) error {
	msgs := make([]*sarama.ProducerMessage, len(message))
	for i, e := range message {
		msgs[i] = s.toMsg(ctx, e)
	}
	return s.writer.SendMessages(msgs)
}

func (s *kafkaSender) toMsg(ctx context.Context, message event.Event) *sarama.ProducerMessage {
	ret := &sarama.ProducerMessage{
		Topic:    s.topic,
		Key:      sarama.StringEncoder(message.Key()),
		Value:    sarama.ByteEncoder(message.Value()),
		Metadata: ctx,
	}
	// push header
	h := message.Header()
	for _, key := range h.Keys() {
		ret.Headers = append(ret.Headers, sarama.RecordHeader{
			Key: []byte(key), Value: []byte(h.Get(key)),
		})
	}

	return ret
}

func (s *kafkaSender) Close() error {
	err := s.writer.Close()
	if err != nil {
		return err
	}
	return nil
}

type kafkaReceiver struct {
	reader  sarama.ConsumerGroup
	topic   string
	group   string
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
	address []string
}

func NewKafkaReceiver(address []string, topic string, group string, cfg *conf.Event_Kafka) (event.Receiver, error) {
	res := &kafkaReceiver{topic: topic, group: group, address: address}
	config := sarama.NewConfig()

	err := patchConf(config, cfg)
	if err != nil {
		return nil, err
	}
	cg, err := sarama.NewConsumerGroup(address, group, config)
	if err != nil {
		return nil, err
	}
	res.reader = cg

	return res, nil

}

func (k *kafkaReceiver) Receive(ctx context.Context, handler event.Handler) error {
	ctx, cancel := context.WithCancel(ctx)
	k.cancel = cancel
	wg := &sync.WaitGroup{}
	wg.Add(1)
	k.wg = wg
	topics := []string{k.topic}
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := k.reader.Consume(ctx, topics, newConsumerGroupHandler(k.group, handler)); err != nil {
				log.Error(err)
				//TODO panic?
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()
	return nil
}

func (k *kafkaReceiver) Close() error {
	k.cancel()
	k.wg.Wait()
	err := k.reader.Close()
	if err != nil {
		return err
	}
	return nil
}

type consumerGroupHandler struct {
	handler event.Handler
	group   string
}

func patchConf(config *sarama.Config, cfg *conf.Event_Kafka) error {
	if cfg != nil {
		if cfg.Version != nil {
			v, err := sarama.ParseKafkaVersion(cfg.Version.Value)
			if err != nil {
				return err
			}
			config.Version = v
		}
	}
	return nil
}

func newConsumerGroupHandler(group string, handler event.Handler) *consumerGroupHandler {
	return &consumerGroupHandler{group: group, handler: handler}
}
func (*consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		//handle msg
		message := event.NewMessage(string(msg.Key), msg.Value)
		for _, header := range msg.Headers {
			message.Header().Set(string(header.Key), string(header.Value))
		}
		err := h.handler.Process(sess.Context(), message)
		if err == nil {
			sess.MarkMessage(msg, "")
		}
	}
	return nil
}
