package kafka

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goxiaoy/go-saas-kit/pkg/event/event"

	"github.com/segmentio/kafka-go"
)

var (
	_ event.Sender   = (*kafkaSender)(nil)
	_ event.Receiver = (*kafkaReceiver)(nil)
)

type kafkaSender struct {
	writer *kafka.Writer
	topic  string
	logger *log.Helper
}

func (s *kafkaSender) Send(ctx context.Context, message event.Event) error {
	err := s.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(message.Key()),
		Value: message.Value(),
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

func NewKafkaSender(address []string, topic string, logger log.Logger) (event.Sender, error) {
	w := &kafka.Writer{
		Topic:    topic,
		Addr:     kafka.TCP(address...),
		Balancer: &kafka.LeastBytes{},
	}
	return &kafkaSender{writer: w, topic: topic, logger: log.NewHelper(logger)}, nil
}

type kafkaReceiver struct {
	reader *kafka.Reader
	topic  string
	logger *log.Helper
}

func (k *kafkaReceiver) Receive(ctx context.Context, handler event.Handler) error {
	go func() {
		for {
			m, err := k.reader.FetchMessage(context.Background())
			if err != nil {
				break
			}
			err = handler(context.Background(), event.NewMessage(string(m.Key), m.Value))
			if err != nil {
				//TODO error handling
				k.logger.Error(fmt.Sprintf("message handling exception: %v", err))
				continue
			}
			if err := k.reader.CommitMessages(ctx, m); err != nil {
				k.logger.Error(fmt.Sprintf("failed to commit messages: %v", err))
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

func NewKafkaReceiver(address []string, topic string, group string, logger log.Logger) (event.Receiver, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  address,
		GroupID:  group,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &kafkaReceiver{reader: r, topic: topic, logger: log.NewHelper(logger)}, nil
}
