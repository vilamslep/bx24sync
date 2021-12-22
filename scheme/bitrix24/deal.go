package bitrix24

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/vi-la-muerto/bx24sync/scheme/bitrix24/converter"
	"github.com/vi-la-muerto/bx24sync/scheme/sql"
)

const (
	findDeal   = "crm.deal.list"
	addDeal    = "crm.deal.add"
	updateDeal = "crm.deal.update"
)

type Deal struct {
	Id          string  `json:"ORIGIN_ID"`
	Name        string  `json:"NAME"`
	User        int     `json:"ASSIGNED_BY_ID"`
	Date        string  `json:"BEGINDATE"`
	Category    int     `json:"CATEGORY_ID"`
	ContactData Contact `json:"ContactData"`
	ContactId   int     `json:"CONTACT_ID"`
	Sum         float32 `json:"OPPORTUNITY"`
	Stage       string  `json:"STAGE_ID"`
	UsersFields []UserField
}

type UserField struct {
	Id    string
	Value string
}

func NewDealFromJson(raw []byte) (deal Deal, err error) {
	err = json.Unmarshal(raw, &deal)
	return deal, err
}

func GetDealFromRawAsReception(reader io.Reader) (data [][]byte, err error) {
	content, err := io.ReadAll(reader)

	if err != nil {
		return data, err
	}

	reseptions := make([]sql.Reception, 0)

	err = json.Unmarshal(content, &reseptions)

	if err != nil {
		return
	}
	for _, reception := range reseptions {

		deal := newDealFromReception(reception)

		if result, err := deal.Json(); err == nil {
			data = append(data, result)
		} else {
			return data, err
		}
	}

	return data, err
}

func (d Deal) Find(restUrl string) (response BitrixRestResponse, err error) {

	if restUrl[len(restUrl)-1:] == "/" {
		restUrl = restUrl[:len(restUrl)-1]
	}

	url := fmt.Sprintf("%s/%s?filter[ORIGIN_ID]=%s", restUrl, findDeal, d.Id)

	if res, err := execReq("GET", url, nil); err == nil {
		return checkResponse(res)
	} else {
		return response, err
	}
}

func (d Deal) Add(restUrl string) (response BitrixRestResponse, err error) {

	if restUrl[len(restUrl)-1:] == "/" {
		restUrl = restUrl[:len(restUrl)-1]
	}

	url := fmt.Sprintf("%s/%s", restUrl, addDeal)

	err = d.checkContact(restUrl)
	if err != nil {
		return
	}

	rd, err := prepareDeal(d)
	if err != nil {
		return
	}

	if res, err := execReq("POST", url, rd); err == nil {
		return checkResponse(res)
	} else {
		return response, err
	}
}

func (d *Deal) checkContact(url string) error {
	response, err := d.ContactData.Find(url)
	if err != nil {
		return err
	}

	if response.Total == 0 {
		if response, err := d.ContactData.Add(url); err == nil {
			id := response.Result[0].ID
			d.ContactId = converter.String(id).Int()
		} else {
			return err
		}
	} else {
		id := response.Result[0].ID
		d.ContactId = converter.String(id).Int()
	}

	return nil
}

func prepareDeal(d Deal) (rd io.Reader, err error) {

	data := make(map[string]map[string]string)

	deal := make(map[string]string)

	deal["ORIGIN_ID"] = d.Id
	deal["NAME"] = d.Name

	if d.User != 0 {
		deal["ASSIGNED_BY_ID"] = fmt.Sprint(d.User)
	}

	deal["BEGINDATE"] = d.Date
	deal["CATEGORY_ID"] = fmt.Sprint(d.Category)
	deal["CONTACT_ID"] = fmt.Sprint(d.ContactId)

	for _, v := range d.UsersFields {
		deal[v.Id] = v.Value
	}

	data["fields"] = deal

	if content, err := json.Marshal(data); err == nil {
		return bytes.NewReader(content), err
	} else {
		return rd, err
	}
}

func (d Deal) Update(restUrl string, id string) (response BitrixRestResponse, err error) {
	if restUrl[len(restUrl)-1:] == "/" {
		restUrl = restUrl[:len(restUrl)-1]
	}

	url := fmt.Sprintf("%s/%s?id=%s", restUrl, updateDeal, id)

	err = d.checkContact(restUrl)
	if err != nil {
		return
	}

	rd, err := prepareDeal(d)
	if err != nil {
		return
	}

	if res, err := execReq("POST", url, rd); err == nil {
		return checkResponse(res)
	} else {
		return response, err
	}
}

func newDealFromReception(reception sql.Reception) Deal {

	d := Deal{}

	d.Id = reception.OriginId
	d.Name = reception.Name
	d.User = converter.String(reception.UserId).Int()
	d.Date = reception.Date
	d.Category = 5
	d.ContactData = newContactFromClient(reception.Client)
	d.Stage = "C5:NEW"

	for _, addFld := range reception.AdditionnalFields {
		var key, value string
		key = addFld.Key

		if converter.String(addFld.Value).IsBinaryBool() {
			value = "0"
			if converter.String(addFld.Value).BinaryTrue() {
				value = "1"
			}
		} else {
			value = addFld.Value
		}

		d.UsersFields = append(d.UsersFields, UserField{
			Id:    key,
			Value: value,
		})
	}
	// UsersFields []UserField

	// c.transforName(client.Name)

	// c.DiscountMedicalThings = converter.String(client.DiscountMedicalThings).Uint8()
	// c.DiscountRayban = converter.String(client.DiscountRayban).Uint8()
	// c.DiscountClinicService = converter.String(client.DiscountClinicService).Uint8()

	// c.Gender = converter.String(client.Gender).GenderCode()
	// c.ConnectionWay = converter.String(client.ConnectionWay).ConnectionWay()

	// c.ThereIsContract = converter.String(client.ThereIsContract).BinaryTrue()
	// c.IsOfflineClient = converter.String(client.IsOfflineClient).BinaryTrue()
	// c.IsInternetClient = converter.String(client.IsInternetClient).BinaryTrue()
	// c.SendAds = converter.String(client.SendAds).BinaryTrue()
	// c.IsClinicClient = converter.String(client.IsClinicClient).BinaryTrue()

	// c.Phone = converter.String(client.Phone).ContactDataSlice(";", "PHONE", "WORK")
	// c.Email = converter.String(client.Email).ContactDataSlice(";", "PHONE", "EMAIL")

	return d
}

func (d *Deal) Json() ([]byte, error) {
	return json.Marshal(d)
}
