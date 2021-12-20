package consumer

// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/vi-la-muerto/bx24sync/app"
// 	bx24 "github.com/vi-la-muerto/bx24sync/bitrix24"
// )

// type newFromMap func(map[string]string) bx24.Entity

// func RunToCrm() {

// 	if err := runScanner(); err != nil {
// 		log.Fatalln(err)
// 	}
// }

// func runScanner() error {
// 	scanner := app.NewKafkaScanner()

// 	for scanner.Scan() {
// 		msg := scanner.Message()
// 		sendMessageToGenerator(msg)
// 	}
// 	return scanner.Err()
// }

// //TODO need to make up where save errors
// func commitError(msg app.Message, err error) {
// 	content := fmt.Sprintf("%s; Error: %s", msg.String(), err.Error())
// 	if ioutil.WriteFile("errors.txt", []byte(content), os.ModeAppend); err != nil {
// 		log.Println(err)
// 	}
// }

// func sendMessageToGenerator(msg app.Message) {

// 	var creating newFromMap
// 	var url string

// 	key := string(msg.Key)

// 	switch key {
// 	case "client":
// 		url = "http://95.78.174.89:25473/client"

// 		creating = bx24.NewContactFromMap

// 	default:
// 		err := fmt.Errorf("not define method for key '%s'", string(msg.Key))
// 		commitError(msg, err)
// 		return
// 	}

// 	reader := bytes.NewReader(msg.Value)

// 	req, err := http.NewRequest("POST", url, reader)

// 	if err != nil {
// 		commitError(msg, err)
// 		return
// 	}

// 	client := http.Client{Timeout: time.Second * 300}

// 	if response, err := client.Do(req); err == nil || response.StatusCode != http.StatusOK {

// 		if err := commitNewMessage(response.Body, creating, key); err != nil {
// 			commitError(msg, err)
// 		}
// 	} else {
// 		commitError(msg, err)
// 	}
// }

// func commitNewMessage(r io.Reader, creating newFromMap, key string) (err error) {

// 	data, err := convertDataForCrm(r, creating)

// 	if err != nil {
// 		return fmt.Errorf("converting for crm failed: %s", err.Error())
// 	}

// 	if err := sendMessageToRegistrar(data, key); err != nil {
// 		return fmt.Errorf("sending message to crm bus failed: %s", err.Error())
// 	}

// 	return nil
// }

// func sendMessageToCRM(content [][]byte, key string) error {

// 	url := fmt.Sprintf("%s/%s", "http://95.78.174.89:25473", key)

// 	for _, data := range content {
// 		rdr := bytes.NewReader(data)

// 		req, err := http.NewRequest("POST", url, rdr)

// 		if err != nil {
// 			return err
// 		}

// 		client := http.Client{Timeout: time.Second * 300}

// 		if response, err := client.Do(req); err == nil && response.StatusCode != http.StatusOK {
// 			return fmt.Errorf("status code isn't expected. Code %d", response.StatusCode)
// 		} else {
// 			return err
// 		}
// 	}

// 	return nil
// }
