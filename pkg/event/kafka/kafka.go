package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
)

var (
	_ event.Sender   = (*kafkaSender)(nil)
	_ event.Receiver = (*kafkaReceiver)(nil)
)

type kafkaSender struct {
	writer sarama.SyncProducer
	topic  string
	logger *log.Helper
}

func NewKafkaSender(address []string, topic string, logger log.Logger) (event.Sender, func(), error) {
	conf := sarama.NewConfig()
	conf.Producer.Interceptors = []sarama.ProducerInterceptor{NewOTelInterceptor(KindProducer, address)}
	w, err := sarama.NewSyncProducer(address, conf)
	if err != nil {
		return nil, func() {

		}, err
	}
	res := &kafkaSender{writer: w, topic: topic, logger: log.NewHelper(log.With(logger, "module", "kafka.kafkaSender"))}
	return res, func() {
		res.Close()
	}, nil
}

func (s *kafkaSender) Send(ctx context.Context, message event.Event) error {
	_, _, err := s.writer.SendMessage(&sarama.ProducerMessage{
		Topic:    s.topic,
		Key:      sarama.StringEncoder(message.Key()),
		Value:    sarama.ByteEncoder(message.Value()),
		Metadata: ctx,
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *kafkaSender) Close() error {
	err := s.writer.Close()
	if err != nil {
		return err
	}
	return nil
}

type kafkaReceiver struct {
	reader sarama.ConsumerGroup
	topic  string
	group  string
	logger *log.Helper
	oi     *OTelInterceptor
}

func NewKafkaReceiver(address []string, topic string, group string, logger log.Logger) (event.Receiver, func(), error) {
	res := &kafkaReceiver{topic: topic, group: group, logger: log.NewHelper(log.With(logger, "module", "kafka.kafkaReceiver"))}
	res.oi = NewOTelInterceptor(KindConsumer, address)
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	cg, err := sarama.NewConsumerGroup(address, group, config)
	if err != nil {
		return nil, func() {

		}, err
	}
	res.reader = cg

	return res, func() {
		res.Close()
	}, nil

}

func (k *kafkaReceiver) Receive(ctx context.Context, handler event.Handler) error {
	topics := []string{k.topic}
	go func() {
		for {
			err := k.reader.Consume(ctx, topics, newConsumerGroupHandler(k.group, handler, k.oi, k.logger))
			if err != nil {
				k.logger.Error(err)
			}
		}
	}()
	return nil
}

func (k *kafkaReceiver) Close() error {
	err := k.reader.Close()
	if err != nil {
		return err
	}
	return nil
}

type consumerGroupHandler struct {
	handler event.Handler
	logger  *log.Helper
	oi      *OTelInterceptor
	group   string
}

func newConsumerGroupHandler(group string, handler event.Handler, oi *OTelInterceptor, logger *log.Helper) *consumerGroupHandler {
	return &consumerGroupHandler{group: group, handler: handler, oi: oi, logger: logger}
}
func (*consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		ctx := context.Background()
		//start span
		ctx, span := h.oi.StartConsumerSpan(ctx, h.group, msg)
		h.logger.Debugf("Message topic:%q key:%s partition:%d offset:%d", msg.Topic, string(msg.Key), msg.Partition, msg.Offset)
		//handle msg
		err := h.handler(sess.Context(), event.NewMessage(string(msg.Key), msg.Value))
		if err != nil {
			h.logger.Error(err)
		} else {
			sess.MarkMessage(msg, "")
		}
		span.End()
	}
	return nil
}
