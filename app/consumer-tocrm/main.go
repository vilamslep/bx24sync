package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"context"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"github.com/vi-la-muerto/bx24-service/scheme/bitrix24"
)

//env
const (
	bx24RestUrl = "BITRIX24_REST_URL"
	kBrokers    = "KAFKA_BROKERS"
	ktopic      = "KAFKA_TOPIC"
	kGroupId    = "KAFKA_CONSUMER_GROUP"
)

func main() {

	var restUrl, brokers, topic, groupId string

	restUrl = getEnvWithFallback(bx24RestUrl, "https://domain.bitrix24.ru/rest/")
	brokers = getEnvWithFallback(kBrokers, "bootstrap-server")
	topic = getEnvWithFallback(ktopic, "tocrm")
	groupId = getEnvWithFallback(kGroupId, "crm")

	//new chanel and exec subscription for handing
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	needToClose := make(chan bool, 1)
	exit := make(chan bool, 1)

	go Run(
		getKafkaReader(brokers, topic, groupId),
		restUrl,
		needToClose,
		exit)

	wait(done, needToClose, exit)
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

func Run(reader *kafka.Reader, url string, needToClose chan bool, exit chan bool) {
	defer reader.Close()

	for {
		select {
		case <-needToClose:
			exit <- true
			return
		default:
			if m, err := reader.ReadMessage(context.Background()); err == nil {
				if crmErr := sendContactToCrm(&m, url); crmErr != nil {
					log.Errorf("sending to crm", crmErr.Error())
					time.Sleep(time.Minute)
				}
			} else {
				log.Errorf("reading from topic: %s\n", err.Error())
				time.Sleep(time.Minute)
				continue
			}
		}
	}
}

func sendContactToCrm(message *kafka.Message, url string) error {

	client := bitrix24.Contact{}
	if decErr := json.Unmarshal(message.Value, &client); decErr != nil {
		return decErr
	}

	oridinId := client.Id

	if id, err := findContact(oridinId, url); err == nil {
		if id == "" {
			return addContact(url, client)
		} else {
			return updateContact(id, url, client)
		}
	}

	return nil
}

func findContact(originId string, rootUrl string) (string, error) {

	workUrl := fmt.Sprintf("%s%s?filter[ORIGIN_ID]=%s", rootUrl, "crm.contact.list", originId)

	if response, err := createAndExecuteRequest(makeGETRequest, workUrl, nil); err == nil {

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return "", err
		}

		if response.StatusCode != http.StatusOK {
			return "", fmt.Errorf("bad requests: %s", string(body))
		}

		result := bitrix24.BitrixRestResponse{}

		decErr := json.Unmarshal(body, &result)

		if decErr != nil {
			return "", err
		}

		if result.Total == 0 {
			return "", nil
		} else {
			return result.Result[0].ID, nil
		}
	} else {
		return "", err
	}
}

func addContact(rootUrl string, client bitrix24.Contact) error {

	workUrl := fmt.Sprintf("%s%s", rootUrl, "crm.contact.add")

	data := make(map[string]bitrix24.Contact)
	data["fields"] = client

	content, encErr := json.Marshal(data)

	breader := bytes.NewReader(content)

	if encErr != nil {
		return encErr
	}

	if response, err := createAndExecuteRequest(makePOSTRequest, workUrl, breader); err == nil {
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			body, err := ioutil.ReadAll(response.Body)

			if err != nil {
				return err
			}
			return fmt.Errorf(fmt.Sprintf("%s: %s", "bad request", string(body)))
		}
	} else {
		return err
	}
	return nil
}

func updateContact(id string, rootUrl string, client bitrix24.Contact) error {

	workUrl := fmt.Sprintf("%s%s?id=%s", rootUrl, "crm.contact.update", id)

	data := make(map[string]bitrix24.Contact)
	data["fields"] = client

	content, encErr := json.Marshal(data)

	breader := bytes.NewReader(content)

	if encErr != nil {
		return encErr
	}

	if response, err := createAndExecuteRequest(makePOSTRequest, workUrl, breader); err == nil {
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			body, err := ioutil.ReadAll(response.Body)

			if err != nil {
				return err
			}
			return fmt.Errorf(fmt.Sprintf("%s: %s", "bad request", string(body)))
		}
	} else {
		return err
	}
	return nil
}

func createAndExecuteRequest(getReq func(string, io.Reader) (*http.Request, error), url string, body io.Reader) (*http.Response, error) {

	request, err := getReq(url, body)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Minute * 1}

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func makeGETRequest(url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest("GET", url, body)
}

func makePOSTRequest(url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest("POST", url, body)
}

func wait(done chan os.Signal, needToClose chan bool, exit chan bool) {
	<-done
	needToClose <- true
	<-exit
}
