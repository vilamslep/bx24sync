package changes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	kafka "github.com/segmentio/kafka-go"
	bx24 "github.com/vi-la-muerto/bx24sync"
)

const (
	kHostKey  = "KAFKA_HOST"
	kPortKey  = "KAFKA_PORT"
	kTopicKey = "KAFKA_TOPIC"
	hsHostKey = "HTTP_HOST"
	hsPortKey = "HTTP_PORT"

	kHostStd  = "kafka"
	kPortStd  = 9092
	kTopicStd = "changes"
	hsHostStd = "localhost"
	hsPortStd = 8082
)

func Run() (err error) {
	config := getConfigFromEnv()
	router := bx24.NewRouter(os.Stdout, os.Stderr)

	writer := getKafkaWriter(config.Broker.String(), config.Topic)

	router.AddMethod(bx24.HttpMethod{
		Path:       "/client",
		Handler:    handlerClient(writer),
		CheckInput: handlerCheckInput,
	})

	server := &http.Server{
		Addr:    config.Http.String(),
		Handler: router,
	}

	func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Can't start to listener: %s\n", err)
		}
	}()

	return err
}

func getConfigFromEnv() bx24.RegistrarConfig {
	return bx24.RegistrarConfig{
		Http: bx24.Socket{
			Host: getEnvWithFallback(hsHostKey, hsHostStd),
			Port: stringToInt(os.Getenv(hsPortKey), hsPortStd),
		},
		ProducerConfig: bx24.ProducerConfig{
			Broker: bx24.Socket{
				Host: getEnvWithFallback(kHostKey, kHostStd),
				Port: stringToInt(os.Getenv(kPortKey), kPortStd),
			},
			Topic: getEnvWithFallback(kTopicKey, kTopicStd),
		},
	}
}

func getEnvWithFallback(env string, fallback string) string {
	val := os.Getenv(env)
	if len(val) == 0 {
		return fallback
	}
	return val
}

func stringToInt(val string, fallback int) int {
	if res, err := strconv.Atoi(val); err != nil {
		return fallback
	} else {
		return res
	}
}

func getKafkaWriter(kafkaURL string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

//handler
func handlerClient(writer *kafka.Writer) func(http.ResponseWriter, *http.Request) error {
	return func(wrt http.ResponseWriter, req *http.Request) error {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			return err
		}
		return fmt.Errorf(string(body))
	}
}

//checkinput
func handlerCheckInput(reader io.Reader) (bool, error) {
	body, err := io.ReadAll(reader)

	if err != nil {
		return false, err
	}

	content := strings.ReplaceAll(string(body), "\n", "")

	regStr := `^{"#",+[[:xdigit:]]{8}(-[[:xdigit:]]{4}){3}-[[:xdigit:]]{12},[\d]{1,6}:[[:xdigit:]]{32}}$`

	if matched, err := regexp.MatchString(regStr, content); err != nil {
		return false, err
	} else {
		return matched, nil
	}
}
