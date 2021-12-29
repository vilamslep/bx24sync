package sql

import "encoding/json"

type Shipment struct {
	Id                 string `json:"originId"`
	Date               string `json:"docDate"`
	Name               string `json:"name"`
	Client             string `json:"client"`
	Sum                string `json:"docSum"`
	InternetOrderStage string `json:"internetOrderStage"`
	Dpp                string `json:"dpp"`
	DppOd              string `json:"dppOD"`
	DppOs              string `json:"dppOS"`
	PickUpPoint        string `json:"pickUpPoint"`
	InternetOrder      string `json:"intenterOrder"`
	SentSms            string `json:"sentSms"`
	OrderType          string `json:"orderType"`
	DeliverySum        string `json:"deliverySum"`
	DeliveryWay        string `json:"deliveryWay"`
	DeliveryAddress    string `json:"deliveryAddress"`
	DeliveryArea       string `json:"deliveryArea"`
	DeliveryTimeFrom   string `json:"deliveryTimeFrom"`
	DeliveryTimeTo     string `json:"deliveryTimeTo"`
	WantedDateShipment string `json:"wantedDateShipment"`
	ExtraInfo          string `json:"extraInfo"`
	Comment            string `json:"comment"`
	Agreement          string `json:"agreement"`
	Stock              string `json:"stock"`
	Doctor             string `json:"doctor"`
	Prepaid            string `json:"prepaid"`
	Prepayment         string `json:"prepayment"`
	Credit             string `json:"credit"`
	Segment			  []Segment `json:"segment"` 
}

func ConvertToShipments(scheme []Field, data []map[string]string) ([]Shipment, error) {
	res := make([]Shipment, 0, len(data))

	for _, v := range data {
		if s, err := (Shipment{}).Convert(scheme, v); err != nil {
			res = append(res, s)
		} else {
			return nil, err
		}
	}
	return res, nil
}

func (s Shipment) Convert(scheme []Field, data map[string]string) (c Shipment, err error) {
	if err := checkByScheme(data, scheme); err != nil {
		return Shipment{}, err
	}

	content, err := json.Marshal(&data)
	if err != nil {
		return Shipment{}, err
	}

	err = json.Unmarshal(content, &s)

	return c, err
}

