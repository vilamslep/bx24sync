package main

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	kafka "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

//env
const (
	kBrokers  = "KAFKA_BROKERS"
	kGroupID  = "KAFKA_CONSUMER_GROUP_ID"
	kTopic    = "KAFKA_TOPIC"
	dEndpoint = "DATABASE_ENDPOINT"
	rEndpoint = "REGISTRAR_ENDPOINT"
)

func main() {

	var brokers, groupId, topic, dbEndpoint, registrarEndpoint string

	brokers = getEnvWithFallback(kBrokers, "bootstrap-server")
	groupId = getEnvWithFallback(kGroupID, "preparing")
	topic = getEnvWithFallback(kTopic, "changes")
	dbEndpoint = getEnvWithFallback(dEndpoint, "http://host/client")
	registrarEndpoint = getEnvWithFallback(rEndpoint, "http://registrar/to-crm")

	//new chanel and exec subscription for handing
	done := make(chan os.Signal, 1)

	needToClose := make(chan bool, 1)
	exit := make(chan bool, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	reader := getKafkaReader(brokers, topic, groupId)

	go Run(reader, dbEndpoint, registrarEndpoint, needToClose, exit)

	//wait signal
	<-done

	needToClose <- true

	<-exit

}

func getEnvWithFallback(env string, fallback string) string {
	val := os.Getenv(env)
	if len(val) == 0 {
		return fallback
	}
	return val
}

func getKafkaReader(brokers string, topic string, groupId string) *kafka.Reader {
	slBrokers := strings.Split(brokers, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: slBrokers,
		GroupID: groupId,
		Topic:   topic,
	})
}

func Run(
	reader *kafka.Reader,
	dbEndpoint string,
	registrarEndpoint string,
	needToClose chan bool,
	exit chan bool) {

	log.Info("Start service")

	defer reader.Close()
	for {
		select {
		case <-needToClose:
			exit <- true
			return
		default:
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Errorf("reading from topic: %s\n", err.Error())
				time.Sleep(time.Minute)
				continue
			}
			if data, err := getDataFromDatabase(&m, dbEndpoint); err == nil {

				reqErr := sendToRegistrar(data, registrarEndpoint)

				if reqErr != nil {
					log.Errorf("sending to registrator:%s\n", reqErr.Error())
					time.Sleep(time.Minute)
				}

			} else {
				log.Errorf("getting from dbms host", err.Error())
				time.Sleep(time.Minute)
			}
		}

	}
}

func getDataFromDatabase(message *kafka.Message, url string) ([]byte, error) {

	data := []byte(message.Value)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/text")

	client := &http.Client{Timeout: time.Second * 90}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	log.Info(string(body))

	return body, nil
}

func sendToRegistrar(data []byte, url string) error {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/text")

	client := &http.Client{Timeout: time.Second * 90}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(string(body))
	}

	return nil
}
