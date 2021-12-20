package sql

import (
	"encoding/json"
	"fmt"
	"io"
)

type Client struct {
	OriginId                  string `json:"originId"`
	Name                      string `json:"name"`
	Birthday                  string `json:"birthday"`
	Gender                    string `json:"gender"`
	IsClient                  string `json:"isClient"`
	IsSuppler                 string `json:"isSuppler"`
	OtherRelation             string `json:"otherRelation"`
	IsRetireeOrDisabledPerson string `json:"isRetireeOrDisabledPerson"`
	ConnectionWay             string `json:"connectionWay"`
	ThereIsContract           string `json:"thereIsContract"`
	SendAds                   string `json:"sendAds"`
	IsInternetClient          string `json:"isInternetClient"`
	IsOfflineClient           string `json:"isOfflineClient"`
	IsClinicClient            string `json:"isClinicClient"`
	DiscountClinicService     string `json:"discountClinicService"`
	DiscountMedicalThings     string `json:"discountMedicalThings"`
	DiscountRayban            string `json:"discountRayban"`
	Phone                     string `json:"phone"`
	Email                     string `json:"email"`
}

type field struct {
	Key      string
	Required bool
}

func ConvertThroughClient(data []map[string]string, scheme []field) (b []byte, err error) {

	res := make([]Client, 0)

	for _, v := range data {
		if c, err := convertToClient(v, scheme); err == nil {
			res = append(res, c)
		} else {
			return nil, err
		}
	}

	return json.Marshal(res)
}

func CreateClientScheme(reader io.Reader) (scheme []field, err error) {
	if content, err := io.ReadAll(reader); err == nil {
		if err := json.Unmarshal(content, &scheme); err != nil {
			return scheme, err
		}
	} else {
		return scheme, err
	}
	return scheme, err
}

func convertToClient(data map[string]string, scheme []field) (c Client, err error) {
	if err := checkByScheme(data, scheme); err != nil {
		return Client{}, err
	}

	content, err := json.Marshal(&data)
	if err != nil {
		return Client{}, err
	}

	err = json.Unmarshal(content, &c)

	return c, err
}

func checkByScheme(data map[string]string, scheme []field) (err error) {
	for _, f := range scheme {
		if _, ok := data[f.Key]; !ok {
			if f.Required {
				err = fmt.Errorf("not found item by key %s", f.Key)
				break
			}
		}
	}
	return err
}

