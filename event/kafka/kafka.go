package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/event"
	"github.com/goava/di"
	"strings"
	"sync"
)

var (
	_ event.Producer = (*Producer)(nil)
	_ event.Consumer = (*Consumer)(nil)
)

func init() {
	event.RegisterConsumer("kafka", func(ctx context.Context, cfg *event.Config, _ *di.Container) (event.Consumer, error) {
		var addr []string
		if cfg.Addr != "" {
			addr = strings.Split(cfg.Addr, ";")
		} else {
			addr = []string{"localhost:9092"}
		}
		return NewConsumer(addr, cfg.Topic, cfg.Group, cfg.Kafka)
	})

	event.RegisterProducer("kafka", func(cfg *event.Config, _ *di.Container) (*event.ProducerMux, error) {
		var addr []string
		if cfg.Addr != "" {
			addr = strings.Split(cfg.Addr, ";")
		} else {
			addr = []string{"localhost:9092"}
		}
		sender, err := NewProducer(addr, cfg.Topic, cfg.Kafka)
		if err != nil {
			return nil, err
		}
		res := event.NewProducer(sender)
		return res, nil
	})
}

type Producer struct {
	sarama.SyncProducer
	topic   string
	address []string
}

func NewProducer(address []string, topic string, cfg *event.Config_Kafka) (*Producer, error) {
	conf := sarama.NewConfig()

	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true

	err := patchConf(conf, cfg)
	if err != nil {
		return nil, err
	}
	w, err := sarama.NewSyncProducer(address, conf)
	if err != nil {
		return nil, err
	}
	s := &Producer{SyncProducer: w, topic: topic, address: address}
	return s, nil
}

func (s *Producer) Send(ctx context.Context, message event.Event) error {
	_, _, err := s.SyncProducer.SendMessage(s.toMsg(ctx, message))
	if err != nil {
		return err
	}
	return nil
}

func (s *Producer) BatchSend(ctx context.Context, message []event.Event) error {
	msgs := make([]*sarama.ProducerMessage, len(message))
	for i, e := range message {
		msgs[i] = s.toMsg(ctx, e)
	}
	return s.SyncProducer.SendMessages(msgs)
}

func (s *Producer) toMsg(ctx context.Context, message event.Event) *sarama.ProducerMessage {
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

func (s *Producer) Close() error {
	err := s.SyncProducer.Close()
	if err != nil {
		return err
	}
	return nil
}

type Consumer struct {
	sarama.ConsumerGroup
	topic   string
	group   string
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
	address []string
}

func NewConsumer(address []string, topic string, group string, cfg *event.Config_Kafka) (event.Consumer, error) {
	res := &Consumer{topic: topic, group: group, address: address}
	config := sarama.NewConfig()

	err := patchConf(config, cfg)
	if err != nil {
		return nil, err
	}
	cg, err := sarama.NewConsumerGroup(address, group, config)
	if err != nil {
		return nil, err
	}
	res.ConsumerGroup = cg

	return res, nil

}

func (k *Consumer) Process(ctx context.Context, handler event.ConsumerHandler) error {
	ctx, cancel := context.WithCancel(ctx)
	k.cancel = cancel
	wg := &sync.WaitGroup{}
	wg.Add(1)
	k.wg = wg
	topics := []string{k.topic}
	failureCounter := 0
	go func() {
		defer wg.Done()
		for {
			// `Process` should be called inside an infinite loop, when a
			// server-side rebalance happens, the Consumer session will need to be
			// recreated to get the new claims
			if err := k.ConsumerGroup.Consume(ctx, topics, newConsumerGroupHandler(k.group, handler)); err != nil {
				failureCounter++
				// infinite loop check failed count to prevent tons of log
				if failureCounter > 10 {
					log.Fatal(err)
				} else {
					log.Error(err)
				}
			}
			// check if context was cancelled, signaling that the Consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()
	return nil
}

func (k *Consumer) Close() error {
	k.cancel()
	k.wg.Wait()
	err := k.ConsumerGroup.Close()
	if err != nil {
		return err
	}
	return nil
}

type consumerGroupHandler struct {
	handler event.ConsumerHandler
	group   string
}

func patchConf(config *sarama.Config, cfg *event.Config_Kafka) error {
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

func newConsumerGroupHandler(group string, handler event.ConsumerHandler) *consumerGroupHandler {
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
