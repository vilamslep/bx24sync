package scheme

import (
	"encoding/json"
	"log"
	"time"
)

type Client struct {
	Name string `json:"name"`
	Code string `json:"code"`
	UID	string `json:"uid"`
}

func (s Client) String() string{
	res, err := json.Marshal(s)
	if err != nil {
		log.Panicf("Transform error: %s", err.Error())
	}

	return string(res)	

}

type Order struct {
	Number string    `json:"number"`
	Date   time.Time `json:"date"`
	UID    string    `json:"uid"`
}

type Shipment struct {
	Number string    `json:"number"`
	Date   time.Time `json:"date"`
	UID    string    `json:"uid"`
}

type Reception struct {
	Number string    `json:"number"`
	Date   time.Time `json:"date"`
	UID    string    `json:"uid"`
}
