package change

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	bx24 "github.com/vi-la-muerto/bx24sync"
	scheme "github.com/vi-la-muerto/bx24sync/scheme/bitrix24"
)

type gettingData func(io.Reader) ([][]byte, error)

type commit struct {
	fields  log.Fields
	message string
	level   string
}

func Run() {

	loggerIn := make(chan commit)

	go commitLogMessage(loggerIn, os.Stdout)

	if err := runScanner(loggerIn); err != nil {
		log.Fatalln(err)
	}
}

func runScanner(loggerIn chan commit) error {

	config, err := bx24.NewConsumerConfigFromEnv()

	if err != nil {
		return err
	}

	scanner := bx24.NewKafkaScanner(config)

	msg := bx24.Message {
		 Topic: "topic",
		 Offset: 1,
	}
	// for scanner.Scan() {
	// 	msg := scanner.Message()
 	sendMessageToGenerator(msg, config.GeneratorEndpoint, config.TargetEndpoint, loggerIn)
	// }
	return scanner.Err()
}

func sendMessageToGenerator(msg bx24.Message, generator bx24.Endpoint, target bx24.Endpoint, loggerIn chan commit) {

	var creating gettingData
	var url string

	key := string(msg.Key)
	url = fmt.Sprintf("%s/%s", generator.URL(), key)

	loggerIn <- commit{
		fields:  log.Fields{"msg": msg},
		message: "get new message from bus",
		level:   "info",
	}

	switch key {
	case "client":
		creating = scheme.GetContactsFromRaw
	default:
		err := fmt.Errorf("not define method for key '%s'", string(msg.Key))
		loggerIn <- commit{
			fields:  log.Fields{"msg": msg},
			message: err.Error(),
			level:   "info",
		}
		return
	}

	rd := bytes.NewReader(msg.Value)

	loggerIn <- commit{
		fields:  log.Fields{"msg": msg},
		message: "Getting data from generator",
		level:   "info",
	}

	if response, err := createAndExecRequest("POST", url, rd); err == nil {
		if response.StatusCode != http.StatusOK {
			err := fmt.Errorf("bad response from generator")
			loggerIn <- commit{
				fields:  log.Fields{"msg": msg},
				message: err.Error(),
				level:   "error",
			}
			return
		}
		defer response.Body.Close()

		loggerIn <- commit{
			fields:  log.Fields{"msg": msg},
			message: "Sending data to registrar",
			level:   "info",
		}
		if err := commitNewMessage(response.Body, creating, key, target); err != nil {
			loggerIn <- commit{
				fields:  log.Fields{"msg": msg},
				message: err.Error(),
				level:   "error",
			}
		}
	} else {
		loggerIn <- commit{
			fields:  log.Fields{"msg": msg},
			message: err.Error(),
			level:   "error",
		}
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

func commitLogMessage(loggerIn chan commit, wrt io.Writer) {
	logger := log.New()
	logger.SetLevel(log.DebugLevel)
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	for {
		msg := <-loggerIn

		entry := log.WithFields(msg.fields)
		if msg.level == "error" {
			entry.Error(msg.message)
		} else {
			entry.Info(msg.message)
		}

	}
}
