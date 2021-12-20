package converter

import "strings"

type Gender int
type ConnectionWay int

const (
	Male   Gender = 263
	Female Gender = 264
)

const (
	SMS   ConnectionWay = 515
	Phone ConnectionWay = 516
)

func (s String) GenderCode() Gender {

	result := Gender(0)

	switch Gender(s.Uint8()) {
	case 0:
		result = Male
	case 1:
		result = Female
	}

	return result
}

func (s String) ConnectionWay() ConnectionWay {

	result := ConnectionWay(0)
	switch ConnectionWay(s.Uint8()) {
	case 0:
		result = SMS
	case 1:
		result = Phone
	}

	return result
}

type ContactData struct {
	ValueType string `json:"VALUE_TYPE"`
	Value     string `json:"VALUE"`
	TypeId    string `json:"TYPE_ID"`
}

func (s String) ContactDataSlice(sep string, valType string, idType string ) []ContactData {

	if len(s) == 0 { return nil }

	target := make([]ContactData, 0)

	for _, piece := range strings.Split(string(s), sep) {
		if piece == "" { continue }

		item := ContactData{
			ValueType: valType,
			Value:     piece,
			TypeId:    idType,
		}

		target = append(target, item)
	}

	return target
}




