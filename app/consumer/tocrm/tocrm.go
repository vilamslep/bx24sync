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

	response, err := entity.Find()

	if err != nil {
		commitError(msg, err)
	}

	if response.Result[0].ID != "" {
		_, _ = entity.Add()
	} else {
		_, _ = entity.Update()
	}

}

//TODO need to make up where save errors
func commitError(msg bx24.Message, err error) {
	content := fmt.Sprintf("%s; Error: %s", msg.String(), err.Error())
	log.Println(content)
}
