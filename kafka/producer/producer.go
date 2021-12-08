package producer

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Addr string
	Topic      string
	Partition  int
	*kafka.Writer
}

func (p *Producer)OpenWriter() {

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
		log.Fatal("failed to write messages:", err)
		return err
	}

	return nil

}

func (s *Producer) CloseWriter() error {
	if err := s.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
		return err
	}

	return nil
}
