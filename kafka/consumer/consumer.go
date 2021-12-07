package consumer

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Consumer struct {
	config Config
	reader *kafka.Reader
}

func NewConsumer(config Config) Consumer {

	c := Consumer{}
	c.config = config
	c.reader = getKafkaReader(config)

	return c
}

func getKafkaReader(config Config) *kafka.Reader {
	brokers := strings.Split(config.Brokers, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: config.GroupId,
		Topic:   config.Topic,
	})
}

func (c *Consumer) Run(needToClose chan bool, exit chan bool) {
	defer c.reader.Close()
	for {
		select {
		case <-needToClose:
			c.Close()
			exit <- true

		default:
			m, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				log.Fatalln(err)
			}
			if data, err := c.getDataFromDatabase(&m); err == nil {
				c.sendToRegistrar(data)
			} else {
				log.Fatalln(err)
			}
		}

	}
}

func (c *Consumer) getDataFromDatabase(message *kafka.Message) ([]byte, error) {

	url := c.config.DBMS
	data := []byte(message.Value)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/text")
	req.Header.Set("Host", "localhost")

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("bad response")
	}

	return body, nil
}

func (c *Consumer) sendToRegistrar(data []byte) error {

	url := c.config.CRMDataline

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/text")
	req.Header.Set("Host", "localhost")

	client := &http.Client{Timeout: time.Second * 10}

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

func (c *Consumer) Close() {

	c.reader.Close()

}
