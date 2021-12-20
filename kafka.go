package bx24sync

//wrapper for github.com/segmentio/kafka-go
import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	Topic     string
	Partition int
	Offset    int64
	Key       []byte
	Value     []byte
}

func (m *Message) String() string {
	return fmt.Sprintf("topic:%s; partition:%d; offset:%d;", m.Topic, m.Partition, m.Offset)
}

func convertKafkaMessageToMessage(msg *kafka.Message) Message {
	return Message{
		Topic:     msg.Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
		Key:       msg.Key,
		Value:     msg.Value,
	}
}

func (m Message) convertToKafkaMessage() kafka.Message {
	return kafka.Message{
		Topic:     m.Topic,
		Partition: m.Partition,
		Offset:    m.Offset,
		Key:       m.Key,
		Value:     m.Value,
	}
}

type KafkaScanner struct {
	reader *kafka.Reader
	ctx    context.Context
	msg    Message
	err    error
}

func NewKafkaScanner(config ConsumerConfig) KafkaScanner {
	return KafkaScanner{
		reader: getKafkaReader(config),
		ctx:    context.Background(),
	}
}

func (r *KafkaScanner) setMessage(msg *kafka.Message) {
	r.msg = convertKafkaMessageToMessage(msg)
}

func (r *KafkaScanner) Scan() bool {
	if msg, err := r.reader.ReadMessage(r.ctx); err == nil {
		r.setMessage(&msg)
		return true
	} else {
		r.err = err
		return false
	}
}

func (r *KafkaScanner) Message() Message {
	return r.msg
}

func (r *KafkaScanner) Err() error {
	return r.err
}

func getKafkaReader(config ConsumerConfig) *kafka.Reader {

	brokers := make([]string, len(config.Brokers))

	for i, v := range config.Brokers {
		brokers[i] = v.String()
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		GroupID:   config.GroupId,
		Partition: config.Partition,
		Topic:     config.Topic,
	})
}

type KafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriter(kafkaURL string, topic string) *KafkaWriter {
	return &KafkaWriter{
		writer: getKafkaWriter(kafkaURL, topic),
	}
}

func (w *KafkaWriter) Write(msg Message) (err error) {
	return w.writer.WriteMessages(context.Background(), msg.convertToKafkaMessage())
}

func (w *KafkaWriter) Close() (err error) {
	return err
}

func getKafkaWriter(kafkaURL string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}
