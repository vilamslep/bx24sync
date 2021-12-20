package bitrix24

import (
	"encoding/json"
	"io"
	"strings"
	"unicode"

	"github.com/vi-la-muerto/bx24sync/scheme/bitrix24/converter"
	"github.com/vi-la-muerto/bx24sync/scheme/sql"
)

type Contact struct {
	Id                    string                  `json:"ORIGIN_ID"`
	Name                  string                  `json:"NAME"`
	SecondName            string                  `json:"SECOND_NAME"`
	LastName              string                  `json:"LAST_NAME"`
	Birthday              string                  `json:"BIRTHDATE"`
	Gender                converter.Gender        `json:"UF_CRM_1587720922"`
	ConnectionWay         converter.ConnectionWay `json:"UF_CRM_1612972723496"`
	ThereIsContract       bool                    `json:"UF_CRM_1584964931121"`
	SendAds               bool                    `json:"UF_CRM_1584970848775"`
	IsInternetClient      bool                    `json:"UF_CRM_1585077731269"`
	IsOfflineClient       bool                    `json:"UF_CRM_1585077757363"`
	IsClinicClient        bool                    `json:"UF_CRM_1585077767559"`
	DiscountClinicService uint8                   `json:"UF_CRM_1587720973"`
	DiscountMedicalThings uint8                   `json:"UF_CRM_1587720952"`
	DiscountRayban        uint8                   `json:"UF_CRM_1587720989"`
	Phone                 []converter.ContactData `json:"PHONE"`
	Email                 []converter.ContactData `json:"EMAIL"`
}

func (s *Contact) transforName(val string) {
	//length
	onlyNameAndLastName, fullName := 2, 3

	val = removeNumbers(val)

	pathName := strings.Split(val, " ")

	length := len(pathName)
	if length >= fullName {
		s.LastName, s.Name, s.SecondName = pathName[0], pathName[1], pathName[2]
	} else if length == onlyNameAndLastName {
		s.LastName, s.Name = pathName[0], pathName[1]
	} else {
		s.Name = val
	}
}

func (c Contact) Json() ([]byte, error) {
	return json.Marshal(c)
}

func NewContactFromJson(raw []byte) (contact Contact,err error) {
	err = json.Unmarshal(raw, &contact)
	return contact, err
}

func (c Contact) Find() (response BitrixRestResponse, err error){ return response, err }

func (c Contact) Add() (response BitrixRestResponse, err error){ return response, err }

// func (c Contact) Get() (response BitrixRestResponse, err error){ return response, err }

func (c Contact) Update() (response BitrixRestResponse, err error){ return response, err }

func removeNumbers(val string) (res string) {

	for _, ch := range val {
		if !unicode.IsNumber(ch) {
			res += string(ch)
		}
	}
	return res
}


func GetContactsFromRaw(reader io.Reader) (data [][]byte, err error) {
	content, err := io.ReadAll(reader)

	if err != nil {
		return data, err
	}

	raw := make([][]byte, 0)

	err = json.Unmarshal(content, &raw)

	if err != nil {
		return
	}
	for _, v := range raw {
		client := sql.Client{}
		json.Unmarshal(v, &client)

		if err != nil {
			return data, err
		}

		if result, err := newContactFromClient(client).Json(); err != nil {
			data = append(data, result)
		}
	}

	return data, err
}

func newContactFromClient(client sql.Client) Contact {

	c := Contact{}

	c.Id = client.OriginId
	c.Birthday = client.Birthday

	c.transforName(client.Name)

	c.DiscountMedicalThings = converter.String(client.DiscountMedicalThings).Uint8()
	c.DiscountRayban = converter.String(client.DiscountRayban).Uint8()
	c.DiscountClinicService = converter.String(client.DiscountClinicService).Uint8()

	c.Gender = converter.String(client.Gender).GenderCode()
	c.ConnectionWay = converter.String(client.ConnectionWay).ConnectionWay()

	c.ThereIsContract = converter.String(client.ThereIsContract).BinaryTrue()
	c.IsOfflineClient = converter.String(client.IsOfflineClient).BinaryTrue()
	c.IsInternetClient = converter.String(client.IsInternetClient).BinaryTrue()
	c.SendAds = converter.String(client.SendAds).BinaryTrue()
	c.IsClinicClient = converter.String(client.IsClinicClient).BinaryTrue()

	c.Phone = converter.String(client.Phone).ContactDataSlice(";", "PHONE", "WORK")
	c.Email = converter.String(client.Email).ContactDataSlice(";", "PHONE", "EMAIL")

	return c
}
