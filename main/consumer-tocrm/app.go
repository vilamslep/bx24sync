package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	bx24 "github.com/vi-la-muerto/bx24sync"
	scheme "github.com/vi-la-muerto/bx24sync/scheme/bitrix24"
)

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
			sendToCrm(msg, config.TargetEndpoint)
			<-marker
		}(marker, msg)
	}
	return scanner.Err()
}

func sendToCrm(msg bx24.Message, target bx24.Endpoint) {

	var entity scheme.Entity
	var err error

	key := string(msg.Key)

	commitLogMessage(commit{
		fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
		message: "get new message from bus",
		level:   "info",
	})

	switch key {
	case "client":
		entity, err = scheme.NewContactFromJson(msg.Value)
		if err != nil {
			commitLogMessage(commit{
				fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
				message: err.Error(),
				level:   "error",
			})
			return
		}
	case "order":
		entity, err = scheme.NewDealFromJson(msg.Value)
		if err != nil {
			commitLogMessage(commit{
				fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
				message: err.Error(),
				level:   "error",
			})
			return
		}
	case "shipment":
		entity, err = scheme.NewDealFromJson(msg.Value)
		if err != nil {
			commitLogMessage(commit{
				fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
				message: err.Error(),
				level:   "error",
			})
			return
		}
	case "reception":
		entity, err = scheme.NewDealFromJson(msg.Value)
		if err != nil {
			commitLogMessage(commit{
				fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
				message: err.Error(),
				level:   "error",
			})
			return
		}
	default:
		err := fmt.Errorf("not define method for key '%s'", string(msg.Key))
		commitLogMessage(commit{
			fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
			message: err.Error(),
			level:   "info",
		})
	}

	restUrl := target.URL()

	commitLogMessage(commit{
		fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
		message: "Finiding data in crm",
		level:   "info",
	})
	response, err := entity.Find(restUrl)

	if err != nil {
		commitLogMessage(commit{
			fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
			message: err.Error(),
			level:   "error",
		})
	}

	commitLogMessage(commit{
		fields:  log.Fields{"resposne": response},
		message: "Finiding result",
		level:   "info",
	})

	if response.Total == 0 {
		commitLogMessage(commit{
			fields:  log.Fields{"key": string(msg.Key)},
			message: "Add new entity",
			level:   "info",
		})
		if _, err := entity.Add(restUrl); err != nil {
			commitLogMessage(commit{
				fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
				message: err.Error(),
				level:   "error",
			})
		}
	} else {
		id := response.Result[0].ID
		commitLogMessage(commit{
			fields:  log.Fields{"key": string(msg.Key), "ID": id},
			message: "Update entity",
			level:   "info",
		})
		if _, err := entity.Update(restUrl, id); err != nil {
			commitLogMessage(commit{
				fields:  log.Fields{"key": string(msg.Key), "offset": msg.Offset, "topic": msg.Topic, "value": string(msg.Value)},
				message: err.Error(),
				level:   "error",
			})
		}
	}

}

func commitLogMessage(msg commit) {
	entry := log.WithFields(msg.fields)
	if msg.level == "error" {
		entry.Error(msg.message)
	} else {
		entry.Info(msg.message)
	}
}
