package sql

import (
	"encoding/json"
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

func ConvertToClients(scheme []Field, data []map[string]string) (b []Client, err error) {

	res := make([]Client, 0, len(data))

	for _, v := range data {
		if c, err := convert(scheme, v); err == nil { 
			res = append(res, c) 
		} else { 
			return nil, err 
		}
	}
	return res, err
}



func convert(scheme []Field, data map[string]string) (c Client, err error) {
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



