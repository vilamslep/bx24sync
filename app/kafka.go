package app

import(
	"github.com/segmentio/kafka-go"
)


func GetKafkaWriter(kafkaURL string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func GetKafkaReader(){
	
}