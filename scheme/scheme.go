package scheme

import (
	"time"
)

type Client struct {
	Name string `json:"name"`
	Code string `json:"code"`
	UID	string `json:"uid"`
}

// func (s Client) String() string{
	
// }

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
