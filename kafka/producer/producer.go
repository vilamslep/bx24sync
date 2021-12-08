package producer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Addr      string
	Topic     string
	Partition int
	*kafka.Writer
}

func (p *Producer) OpenWriter() {

	p.Writer = &kafka.Writer{
		Addr:     kafka.TCP(p.Addr),
		Topic:    p.Topic,
		Balancer: &kafka.LeastBytes{},
	}

}

func (s *Producer) WriteMessage(data []byte, key string) error {

	err := s.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:       []byte(key),
			Value:     data,
			Partition: s.Partition,
		},
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *Producer) CloseWriter() error {
	if err := s.Close(); err != nil {
		return err
	}

	return nil
}
