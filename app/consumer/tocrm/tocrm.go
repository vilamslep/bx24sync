package tocrm

import (
	"fmt"
	"log"

	bx24 "github.com/vi-la-muerto/bx24sync"
	scheme "github.com/vi-la-muerto/bx24sync/scheme/bitrix24"
)

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
		sendToCrm(msg, config.TargetEndpoint)
	}
	return scanner.Err()
}

func sendToCrm(msg bx24.Message, target bx24.Endpoint) {
	key := string(msg.Key)

	var entity scheme.Entity
	var err error

	switch key {
	case "client":
		entity, err = scheme.NewContactFromJson(msg.Value)
		if err != nil {
			commitError(msg, err)
			return
		}
	case "order":
	case "shipment":
	case "reception":
	default:
		err := fmt.Errorf("not define method for key '%s'", string(msg.Key))
		commitError(msg, err)
	}

	restUrl := target.URL()

	response, err := entity.Find(restUrl)

	if err != nil {
		commitError(msg, err)
	}

	if response.Total == 0 {
		if _, err := entity.Add(restUrl); err != nil {
			commitError(msg, err)
		}
	} else {
		id := response.Result[0].ID
		if _, err := entity.Update(restUrl, id); err != nil {
			commitError(msg, err)
		}
	}
}

//TODO need to make up where save errors
func commitError(msg bx24.Message, err error) {
	content := fmt.Sprintf("%s; Error: %s", msg.String(), err.Error())
	log.Println(content)
}
