package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	bx24 "github.com/vilamslep/bx24sync"
	scheme "github.com/vilamslep/bx24sync/scheme/bitrix24"
)

type gettingData func(io.Reader) ([][]byte, error)

type commit struct {
	fields  log.Fields
	message string
	level   string
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {

	return runScanner()
}

func runScanner() error {

	config, err := bx24.NewConsumerConfigFromEnv()

	if err != nil {
		return err
	}

	scanner := bx24.NewKafkaScanner(config)

	marker := make(chan struct{}, 20)
	for scanner.Scan() {

		msg := scanner.Message()

		marker <- struct{}{}

		go func(marker chan struct{}, msg bx24.Message) {
			sendMessageToGenerator(msg, config.GeneratorEndpoint, config.TargetEndpoint)
			<-marker
		}(marker, msg)
	}

	return scanner.Err()
}

func sendMessageToGenerator(msg bx24.Message, generator bx24.Endpoint, target bx24.Endpoint) {

	var creating gettingData
	var url string

	key := string(msg.Key)
	url = fmt.Sprintf("%s/%s", generator.URL(), key)

	commitLogMessage(commit{
		fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
		message: "get new message from bus",
		level:   "info",
	})

	switch key {
	case "client":
		creating = scheme.GetContactsFromRaw
	case "reception":
		creating = scheme.GetDealFromRawAsReception
	case "order":
		creating = scheme.GetDealFromRawAsOrder
	case "shipment":
		creating = scheme.GetDealFromRawAsShipment
	default:
		err := fmt.Errorf("not define method for key '%s'", string(msg.Key))
		commitLogMessage(commit{
			fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
			message: err.Error(),
			level:   "info",
		})
		return
	}

	rd := bytes.NewReader(msg.Value)

	commitLogMessage(commit{
		fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
		message: "Getting data from generator",
		level:   "info",
	})

	if response, err := scheme.ExecReq("POST", url, rd); err == nil {
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {

			content, err := io.ReadAll(response.Body)
			if err != nil {
				commitLogMessage(commit{
					fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value), "url": url},
					message: fmt.Errorf("reading response.%s", err.Error()).Error(),
					level:   "error",
				})
			}
			err = fmt.Errorf("bad response from generator. Reponse: %s", string(content))

			commitLogMessage(commit{
				fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value), "url": url},
				message: err.Error(),
				level:   "error",
			})
			return
		}

		commitLogMessage(commit{
			fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
			message: "Sending data to registrar",
			level:   "info",
		})
		if err := commitNewMessage(response.Body, creating, key, target); err != nil {
			commitLogMessage(commit{
				fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
				message: err.Error(),
				level:   "error",
			})
		}
	} else {
		commitLogMessage(commit{
			fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
			message: err.Error(),
			level:   "error",
		})
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

		if response, err := scheme.ExecReq("POST", url, rd); err == nil {
			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				content, err := io.ReadAll(response.Body)

				if err != nil {
					return fmt.Errorf("status code %d. Can't read response. %v", response.StatusCode, err)
				}
				return fmt.Errorf("status code isn't expected. Code %d. Response: %s", response.StatusCode, string(content))
			}

		} else {
			return err
		}
	}
	return nil
}

func commitLogMessage(msg commit) {
	entry := log.WithFields(msg.fields)
	if msg.level == "error" {
		entry.Error(msg.message)
	} else {
		entry.Info(msg.message)
	}
}
