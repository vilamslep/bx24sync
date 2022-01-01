package sql

import "encoding/json"

type Shipment struct {
	Ref        string    `json:"ref"`
	Id         string    `json:"originId"`
	Name       string    `json:"name"`
	Date       string    `json:"docDate"`
	ClientId   string    `json:"client"`
	Client     Client    `json:"clientData,omitempty"`
	Sum        string    `json:"docSum"`
	Department string    `json:"department"`
	Stock      string    `json:"stock"`
	Agreement  string    `json:"agreement"`
	Comment    string    `json:"comment"`
	Doctor     string    `json:"doctor"`
	User       string    `json:"userId"`
	Segments   []Segment `json:"segment"`
}

func ConvertToShipment(scheme []Field, data []map[string]string) ([]Shipment, error) {
	res := make([]Shipment, 0, len(data))

	for _, v := range data {
		if s, err := (Shipment{}).convert(scheme, v); err == nil {
			res = append(res, s)
		} else {
			return nil, err
		}
	}
	return res, nil
}

func (s Shipment) convert(scheme []Field, data map[string]string) (Shipment, error) {

	if err := checkByScheme(data, scheme); err != nil {
		return Shipment{}, err
	}

	content, err := json.Marshal(&data)
	if err != nil {
		return Shipment{}, err
	}

	err = json.Unmarshal(content, &s)

	return s, err
}

func (s *Shipment) LoadSegments(segments []map[string]string) {
	unique := make(map[string][]string)

	for _, m := range segments {
		k, v := m["segment"], m["brand"]
		if _, ok := unique[k]; !ok {
			unique[k] = make([]string, 0)
		}
		unique[k] = append(unique[k], v)
	}

	if s.Segments == nil {
		s.Segments = make([]Segment, 0)
	}

	for k, v := range unique {
		s.Segments = append(s.Segments, Segment{Id: k, Brands: v})
	}
}
