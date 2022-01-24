package bitrix24

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/vilamslep/bx24sync/scheme/bitrix24/converter"
	"github.com/vilamslep/bx24sync/scheme/sql"
)

const (
	findContact   = "crm.contact.list"
	addContact    = "crm.contact.add"
	updateContact = "crm.contact.update"
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
	NotSendAds            bool                    `json:"UF_CRM_1636547489168"`
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

func (c *Contact) Json() ([]byte, error) {
	return json.Marshal(c)
}

func NewContactFromJson(raw []byte) (contact Contact, err error) {
	err = json.Unmarshal(raw, &contact)
	return contact, err
}

func (c Contact) Find(restUrl string) (response BitrixRestResponseFind, err error) {

	if restUrl[len(restUrl)-1:] == "/" {
		restUrl = restUrl[:len(restUrl)-1]
	}

	url := fmt.Sprintf("%s/%s?filter[ORIGIN_ID]=%s", restUrl, findContact, c.Id)

	if res, err := ExecReq("GET", url, nil); err == nil {
		return checkResponseFind(res)
	} else {
		return response, err
	}
}

func (c Contact) Add(restUrl string) (response BitrixRestResponseAdd, err error) {

	if restUrl[len(restUrl)-1:] == "/" {
		restUrl = restUrl[:len(restUrl)-1]
	}

	url := fmt.Sprintf("%s/%s", restUrl, addContact)

	rd, err := c.getReader()
	if err != nil {
		return response, err
	}

	if res, err := ExecReq("POST", url, rd); err == nil {
		return checkResponseAdd(res)
	} else {
		return response, err
	}
}
// for cleaning contact data at first need to get id entity
// array should be such as 
// [ 
// 	"$entityId" => {
// 		"VALUE_TYPE" => "WORK", 
// 		"VALUE" => "", 
// 		"TYPE_ID" => "PHONE" 
// 	}
// ]
func (c Contact) Update(restUrl string, id string) (response BitrixRestResponseUpdate, err error) {
	if restUrl[len(restUrl)-1:] == "/" {
		restUrl = restUrl[:len(restUrl)-1]
	}

	url := fmt.Sprintf("%s/%s?id=%s", restUrl, updateContact, id)

	rd, err := c.getReader()
	if err != nil {
		return response, err
	}

	if res, err := ExecReq("POST", url, rd); err == nil {
		return checkResponseUpdate(res)
	} else {
		return response, err
	}
}

func (c Contact) getReader() (rd io.Reader, err error) {

	data := make(map[string]Contact)

	data["fields"] = c

	if content, err := json.Marshal(data); err == nil {
		return bytes.NewReader(content), err
	} else {
		return rd, err
	}
}

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

	raw := make([]map[string]string, 0)

	err = json.Unmarshal(content, &raw)

	if err != nil {
		return
	}
	for _, v := range raw {
		rwClient, err := json.Marshal(v)
		if err != nil {
			return data, err
		}
		client := sql.Client{}
		json.Unmarshal(rwClient, &client)

		if err != nil {
			return data, err
		}

		contact := newContactFromClient(client)

		if result, err := contact.Json(); err == nil {
			data = append(data, result)
		} else {
			return data, err
		}
	}

	return data, err
}

func newContactFromClient(client sql.Client) Contact {

	c := Contact{}

	offset := 2000

	c.Id = client.OriginId
	c.Birthday = converter.SubtractionYearsOffset(client.Birthday, offset, "02.01.2006")

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

	c.Phone = converter.String(client.Phone).ContactDataSlice(";", "WORK", "PHONE")
	c.Email = converter.String(client.Email).ContactDataSlice(";", "WORK", "EMAIL")

	return c
}
