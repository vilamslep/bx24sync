package consumer

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vi-la-muerto/bx24sync/app"
	bx24 "github.com/vi-la-muerto/bx24sync/scheme/bitrix24"
)

type gettingData func(io.Reader) ([][]byte, error)

func RunPreparing() {

	if err := runScanner(); err != nil {
		log.Fatalln(err)
	}
}

func runScanner() error {
	scanner := app.NewKafkaScanner()

	for scanner.Scan() {
		msg := scanner.Message()
		sendMessageToGenerator(msg)
	}
	return scanner.Err()
}

//TODO need to make up where save errors
func commitError(msg app.Message, err error) {
	content := fmt.Sprintf("%s; Error: %s", msg.String(), err.Error())
	if os.WriteFile("errors.txt", []byte(content), os.ModeAppend); err != nil {
		log.Println(err)
	}
}

func sendMessageToGenerator(msg app.Message) {

	var creating gettingData
	var url string

	key := string(msg.Key)

	switch key {
	case "client":
		url = "http://localhost/client"

		creating = bx24.GetContactsFromRaw

	default:
		err := fmt.Errorf("not define method for key '%s'", string(msg.Key))
		commitError(msg, err)
		return
	}

	reader := bytes.NewReader(msg.Value)

	req, err := http.NewRequest("POST", url, reader)

	if err != nil {
		commitError(msg, err)
		return
	}

	client := http.Client{Timeout: time.Second * 300}

	if response, err := client.Do(req); err == nil || response.StatusCode != http.StatusOK {

		defer response.Body.Close()

		if err := commitNewMessage(response.Body, creating, key); err != nil {
			commitError(msg, err)
		}
	} else {
		commitError(msg, err)
	}
}

func commitNewMessage(r io.Reader, creating gettingData, key string) (err error) {

	data, err := convertDataForCrm(r, creating)

	if err != nil {
		return fmt.Errorf("converting for crm failed: %s", err.Error())
	}

	if err := sendMessageToRegistrar(data, key); err != nil {
		return fmt.Errorf("sending message to crm bus failed: %s", err.Error())
	}

	return nil
}

func convertDataForCrm(r io.Reader, creating gettingData) (data [][]byte, err error) {
	return creating(r)
}

func sendMessageToRegistrar(content [][]byte, key string) error {

	url := fmt.Sprintf("%s/%s", "http://localhost", key)

	for _, data := range content {
		rdr := bytes.NewReader(data)

		req, err := http.NewRequest("POST", url, rdr)

		if err != nil {
			return err
		}

		client := http.Client{Timeout: time.Second * 300}

		response, err := client.Do(req)

		if err != nil {
			return err
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("status code isn't expected. Code %d", response.StatusCode)
		}
	}

	return nil
}
