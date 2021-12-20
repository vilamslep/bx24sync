package change

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	bx24 "github.com/vi-la-muerto/bx24sync"
	scheme "github.com/vi-la-muerto/bx24sync/scheme/bitrix24"
)

type gettingData func(io.Reader) ([][]byte, error)

func Run() {

	if err := runScanner(); err != nil {
		log.Fatalln(err)
	}
}

func runScanner() error {

	config, err := bx24.NewConsumerConfigFromEnv()

	if err != nil {
		return err
	}

	scanner := bx24.NewKafkaScanner(config)

	for scanner.Scan() {
		msg := scanner.Message()
		sendMessageToGenerator(msg, config.GeneratorEndpoint, config.TargetEndpoint)
	}
	return scanner.Err()
}

func sendMessageToGenerator(msg bx24.Message, generator bx24.Endpoint, target bx24.Endpoint) {

	var creating gettingData
	var url string

	key := string(msg.Key)
	url = fmt.Sprintf("%s/%s", generator.URL(), key)

	switch key {
	case "client":
		creating = scheme.GetContactsFromRaw
	default:
		err := fmt.Errorf("not define method for key '%s'", string(msg.Key))
		commitError(msg, err)
		return
	}

	rd := bytes.NewReader(msg.Value)

	if response, err := createAndExecRequest("POST", url, rd); err == nil {
		if response.StatusCode != http.StatusOK {
			err := fmt.Errorf("bad response from generator" )
			commitError(msg, err)
			return
		}
		defer response.Body.Close()

		if err := commitNewMessage(response.Body, creating, key, target); err != nil {
			commitError(msg, err)
		}
	} else {
		commitError(msg, err)
	}
}

func commitNewMessage(r io.Reader, creating gettingData, key string, target bx24.Endpoint) (err error) {

	data, err := convertDataForCrm(r, creating)

	if err != nil {
		return fmt.Errorf("converting for crm failed: %s", err.Error())
	}

	if err := sendMessageToRegistrar(data, key, target); err != nil {
		return fmt.Errorf("sending message to crm bus failed: %s", err.Error())
	}

	return nil
}

func convertDataForCrm(r io.Reader, creating gettingData) (data [][]byte, err error) {
	return creating(r)
}

func sendMessageToRegistrar(content [][]byte, key string, target bx24.Endpoint) error {

	url := fmt.Sprintf("%s/%s", target.URL(), key)

	for _, data := range content {
		rd := bytes.NewReader(data)

		if response, err := createAndExecRequest("POST", url, rd); err == nil {
			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				return fmt.Errorf("status code isn't expected. Code %d", response.StatusCode)
			}

		} else {
			return err
		}
	}
	return nil
}

func createAndExecRequest(method string, url string, rd io.Reader) (*http.Response, error) {

	if req, err := http.NewRequest(method, url, rd); err == nil {
		client := http.Client{Timeout: time.Second * 300}
		return client.Do(req)
	} else {
		return nil, err
	}
}

//TODO need to make up where save errors
func commitError(msg bx24.Message, err error) {
	content := fmt.Sprintf("%s; Error: %s", msg.String(), err.Error())
	log.Println(content)
}
