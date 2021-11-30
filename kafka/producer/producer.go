package producer

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	BrokerAddr string
	Port       int
	Topic      string
	Partition  int
	*kafka.Writer
}

func CreateWriter(p Producer) *kafka.Writer {

	addr := fmt.Sprintf("%s:%d", p.BrokerAddr, p.Port)

	return &kafka.Writer{
		Addr:     kafka.TCP(addr),
		Topic:    p.Topic,
		Balancer: &kafka.LeastBytes{},
	}

}

func (s *Producer) WriteMessage(data string, key string) error {

	err := s.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:       []byte(key),
			Value:     []byte(data),
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
