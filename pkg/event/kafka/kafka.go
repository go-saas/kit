package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"
	"sync"
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

func NewKafkaSender(address []string, topic string, logger log.Logger, cfg *conf.Event_Kafka) (event.Sender, func(), error) {
	conf := sarama.NewConfig()
	conf.Producer.Interceptors = []sarama.ProducerInterceptor{NewOTelInterceptor(KindProducer, address)}
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
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func NewKafkaReceiver(address []string, topic string, group string, logger log.Logger, cfg *conf.Event_Kafka) (event.Receiver, func(), error) {
	res := &kafkaReceiver{topic: topic, group: group, logger: log.NewHelper(log.With(logger, "module", "kafka.kafkaReceiver"))}
	res.oi = NewOTelInterceptor(KindConsumer, address)
	config := sarama.NewConfig()
	var cleanup = func() {

	}
	err := patchConf(config, cfg)
	if err != nil {
		return nil, cleanup, err
	}
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
			if err := k.reader.Consume(ctx, topics, newConsumerGroupHandler(k.group, handler, k.oi, k.logger)); err != nil {
				k.logger.Error(err)
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
	logger  *log.Helper
	oi      *OTelInterceptor
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
		h.oi.EndConsumerSpan(ctx, span, err)
	}
	return nil
}
