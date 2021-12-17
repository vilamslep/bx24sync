package consumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vi-la-muerto/bx24sync/app"
	bx24 "github.com/vi-la-muerto/bx24sync/bitrix24"
)

func Run()  {

	if err := runScanner(); err != nil {
		log.Fatalln(err)
	}



}

func runScanner() error {
	scanner := app.NewKafkaScanner()

	for scanner.Scan() {
		msg := scanner.Message
		go sendMessageToGenerator(msg)
	}
	return scanner.GetError()
}
//TODO need to make up where save errors
func commitError(msg app.Message, err error) {
	content := fmt.Sprintf("%s; Error: %s", msg.String(), err.Error())
	if ioutil.WriteFile("errors.txt", []byte(content), os.ModeAppend); err != nil {
		log.Println(err)
	}
}

func sendMessageToGenerator(msg app.Message) {
	var entity bx24.Entity
	var url string
	switch string(msg.Key) {
	case  "client":
		url = fmt.Sprintf("%s/%s", "http://192.168.2.238:8095", msg.Key)
		entity = bx24.Contact{}
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

	client := http.Client{ Timeout: time.Second * 300}

	if response, err := client.Do(req); err == nil || response.StatusCode != http.StatusOK {

		if err := commitNewMessage(response.Body, entity); err != nil {
			commitError(msg, err)
		}
	} else {
		commitError(msg, err)
	}
}

func commitNewMessage(r io.Reader, entity bx24.Entity) (err error) {

	data, err := convertDataForCrm(r, entity)

	if err != nil {
		return fmt.Errorf("converting for crm failed: %s", err.Error())
	}

	if err := sendMessageToRegistrar(data); err != nil {
		return fmt.Errorf("sending message to crm bus failed: %s", err.Error())
	}
	
	return nil
}

func convertDataForCrm(r io.Reader, entity bx24.Entity) (data []byte, err error) {

	src := make(map[string]string)
	
	err = json.Unmarshal(data, &src)

	if err != nil {
		return data, err
	}

	entity.LoadFromMap(src);
	
	return entity.Json()
}


func sendMessageToRegistrar(content []byte) error {

	return nil

}