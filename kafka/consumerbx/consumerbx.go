package consumerbx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/vi-la-muerto/bx24-service/http/generator/scheme"
)

type Consumer struct {
	config Config
	reader *kafka.Reader
}

func NewConsumer(config Config) Consumer {

	c := Consumer{}
	c.config = config
	c.reader = getKafkaReader(config)

	return c
}

func getKafkaReader(config Config) *kafka.Reader {
	brokers := strings.Split(config.Brokers, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: config.GroupId,
		Topic:   config.Topic,
	})
}

func (c *Consumer) Run(needToClose chan bool, exit chan bool) {
	defer c.reader.Close()

	for {
		select {
		case <-needToClose:
			c.Close()
			exit <- true

		default:
			m, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				log.Fatalln(err)
			}
			if err := c.sendContactToCrm(&m); err != nil {
				log.Println(err)
			}
		}

	}
}

func (c *Consumer) sendContactToCrm(message *kafka.Message ) error {

	client := scheme.Client{}
	if decErr := json.Unmarshal(message.Value, &client); decErr != nil {
		return decErr
	}
	
	oridinId := client.Id

	rootUrl := c.config.RESTUrl

	if id, err := c.findContact(oridinId, rootUrl); err == nil {
		if id == "" {
			return c.addContact(rootUrl, client)
		} else {
			return c.updateContact(id, rootUrl, client)
		}
	}
	
	return nil
}

func (c *Consumer) findContact(originId string, rootUrl string) (string, error) {

	workUrl := fmt.Sprintf("%s%s?filter[ORIGIN_ID]=%s", rootUrl, "crm.contact.list", originId)

	if response, err := createAndExecuteRequest(makeGETRequest, workUrl, nil); err == nil {
		
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return "", err
		}

		if response.StatusCode != http.StatusOK {
			return "", errors.New("Bad requests")
		}

		result := BitrixRestResponse{}
		
		decErr := json.Unmarshal(body, &result)

		if decErr != nil {
			return "", err
		}

		if result.Total == 0 {
			return "", nil
		} else {
			return result.Result[0].ID, nil
		}

	}else {
		return "", err
	}	
}

func (c *Consumer) addContact(rootUrl string, client scheme.Client) error {

	workUrl := fmt.Sprintf("%s%s", rootUrl, "crm.contact.add")

	data := make(map[string]scheme.Client)
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
			return errors.New(fmt.Sprintf("%s: %s", "bad request", string(body)))
		}
	} else {
		return err
	}
	return nil
}

func (c *Consumer) updateContact(id string, rootUrl string, client scheme.Client) error {

	workUrl := fmt.Sprintf("%s%s?id=%s", rootUrl, "crm.contact.update", id)

	data := make(map[string]scheme.Client)
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
			return errors.New(fmt.Sprintf("%s: %s", "bad request", string(body)))
		}
	} else {
		return err
	}
	return nil
}


func makeGETRequest(url string, body io.Reader) (*http.Request, error){
	return http.NewRequest("GET", url, body)
}

func makePOSTRequest(url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest("POST", url, body)
}

func createAndExecuteRequest( getReq func(string, io.Reader)(*http.Request, error), url string, body io.Reader ) (*http.Response, error) {

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

func (c *Consumer) Close() {
	c.reader.Close()
}