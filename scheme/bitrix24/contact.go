package bitrix24

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/vi-la-muerto/bx24sync/scheme/bitrix24/converter"
	"github.com/vi-la-muerto/bx24sync/scheme/sql"
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

func (c Contact) Find(restUrl string) (response BitrixRestResponse, err error) {

	if restUrl[len(restUrl)-1:] == "/" {
		restUrl = restUrl[:len(restUrl)-1]
	}

	url := fmt.Sprintf("%s/%s?filter[ORIGIN_ID]=%s", restUrl, findContact, c.Id)

	if res, err := execReq("GET", url, nil); err == nil {
		return checkResponse(res)
	} else {
		return response, err
	}
}

func (c Contact) Add(restUrl string) (response BitrixRestResponse, err error) {

	if restUrl[len(restUrl)-1:] == "/" {
		restUrl = restUrl[:len(restUrl)-1]
	}

	url := fmt.Sprintf("%s/%s", restUrl, addContact)

	rd, err := prepareReader(c)
	if err != nil {
		return response, err
	}

	if res, err := execReq("POST", url, rd); err == nil {
		return checkResponse(res)
	} else {
		return response, err
	}
}

func (c Contact) Update(restUrl string, id string) (response BitrixRestResponse, err error) {
	if restUrl[len(restUrl)-1:] == "/" {
		restUrl = restUrl[:len(restUrl)-1]
	}

	url := fmt.Sprintf("%s/%s?id=%s", restUrl, updateContact, id)

	rd, err := prepareReader(c)
	if err != nil {
		return response, err
	}

	if res, err := execReq("POST", url, rd); err == nil {
		return checkResponse(res)
	} else {
		return response, err
	}
}

func prepareReader(c Contact) (rd io.Reader, err error) {

	data := make(map[string]Contact)
	data["fields"] = c

	if content, err := json.Marshal(data); err == nil {
		return bytes.NewReader(content), err
	} else {
		return rd, err
	}
}

func checkResponse(res *http.Response) (response BitrixRestResponse, err error) {

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		if err != nil {
			return response, err
		}
		return response, fmt.Errorf(fmt.Sprintf("%s: %s", "bad request", string(body)))
	} else {
		return fillResponse(body)
	}
}

func fillResponse(bodyRaw []byte) (response BitrixRestResponse, err error) {
	err = json.Unmarshal(bodyRaw, &response)
	return response, err
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

func execReq(method string, url string, rd io.Reader) (*http.Response, error) {

	if req, err := http.NewRequest(method, url, rd); err == nil {
		client := http.Client{Timeout: time.Second * 300}
		return client.Do(req)
	} else {
		return nil, err
	}
}
