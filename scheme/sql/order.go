package sql

import "encoding/json"

type Order struct {
	Ref                string    `json:"ref"`
	Id                 string    `json:"originId"`
	Date               string    `json:"docDate"`
	Name               string    `json:"name"`
	ClientId           string    `json:"client"`
	Client             Client    `json:"clientData,omitempty"`
	Sum                string    `json:"docSum"`
	InternetOrderStage string    `json:"internetOrderStage"`
	Dpp                string    `json:"dpp"`
	DppOD              string    `json:"dppOD"`
	DppOS              string    `json:"dppOS"`
	PickUpPoint        string    `json:"pickUpPoint"`
	InternetOrder      string    `json:"intenterOrder"`
	SentSms            string    `json:"sentSms"`
	OrderType          string    `json:"orderType"`
	DeliverySum        string    `json:"deliverySum"`
	DeliveryWay        string    `json:"deliveryWay"`
	WantedDateShipment string    `json:"wantedDateShipment"`
	ExtraInfo          string    `json:"extraInfo"`
	Comment            string    `json:"comment"`
	Agreement          string    `json:"agreement"`
	Stock              string    `json:"stock"`
	DeliveryAddress    string    `json:"deliveryAddress"`
	DeliveryArea       string    `json:"deliveryArea"`
	DeliveryTimeFrom   string    `json:"deliveryTimeFrom"`
	DeliveryTimeTo     string    `json:"deliveryTimeTo"`
	Doctor             string    `json:"doctor"`
	UserId             string    `json:"userId"`
	ShipmentDate       string    `json:"shipmentDate"`
	Prepaid            string    `json:"prepaid"`
	Prepayment         string    `json:"prepayment"`
	Credit             string    `json:"credit"`
	Department         string    `json:"department"`
	Segments           []Segment `json:"segment"`
}

func ConvertToOrders(scheme []Field, data []map[string]string) ([]Order, error) {
	res := make([]Order, 0, len(data))

	for _, v := range data {
		if s, err := (Order{}).convert(scheme, v); err == nil {
			res = append(res, s)
		} else {
			return nil, err
		}
	}
	return res, nil
}

func (o Order) convert(scheme []Field, data map[string]string) (Order, error) {

	if err := checkByScheme(data, scheme); err != nil {
		return Order{}, err
	}

	content, err := json.Marshal(&data)
	if err != nil {
		return Order{}, err
	}

	err = json.Unmarshal(content, &o)

	return o, err
}

func (o *Order) LoadSegments(segments []map[string]string) {
	//ref segment brand
	unique := make(map[string][]string)

	for _, m := range segments {
		k, v := m["segment"], m["brand"]
		if _, ok := unique[k]; !ok {
			unique[k] = make([]string, 0)
		}
		unique[k] = append(unique[k], v)
	}

	if o.Segments == nil {
		o.Segments = make([]Segment, 0)
	}

	for k, v := range unique {
		o.Segments = append(o.Segments, Segment{Id: k, Brands: v})
	}
}
