package bitrix24

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/vilamslep/bx24sync/scheme/bitrix24/converter"
	"github.com/vilamslep/bx24sync/scheme/sql"
)

const (
	findDeal   = "crm.deal.list"
	addDeal    = "crm.deal.add"
	updateDeal = "crm.deal.update"
)

type Deal struct {
	Id               string  `json:"ORIGIN_ID"`
	Name             string  `json:"NAME"`
	User             int     `json:"ASSIGNED_BY_ID"`
	Date             string  `json:"BEGINDATE"`
	Category         int     `json:"CATEGORY_ID"`
	ContactData      Contact `json:"ContactData"`
	ContactId        int     `json:"CONTACT_ID"`
	Sum              float32 `json:"OPPORTUNITY"`
	Stage            string  `json:"STAGE_ID"`
	Comment          string  `json:"COMMENTS"`
	UsersFields      []UserField
	UserFieldsPlurar []UserFiledPlurarValue
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

func GetDealFromRawAsOrder(reader io.Reader) (data [][]byte, err error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return data, err
	}

	orders := make([]sql.Order, 0)

	err = json.Unmarshal(content, &orders)

	if err != nil {
		return
	}
	for _, order := range orders {

		deal, err := newDealFromOrder(order)
		if err != nil {
			return nil, err
		}

		if result, err := deal.Json(); err == nil {
			data = append(data, result)
		} else {
			return data, err
		}
	}

	return data, err
}

func GetDealFromRawAsShipment(reader io.Reader) (data [][]byte, err error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return data, err
	}

	shipments := make([]sql.Shipment, 0)

	err = json.Unmarshal(content, &shipments)

	if err != nil {
		return
	}
	for _, shipment := range shipments {

		deal, err := newDealFromShipment(shipment)
		if err != nil {
			return nil, err
		}

		if result, err := deal.Json(); err == nil {
			data = append(data, result)
		} else {
			return data, err
		}
	}

	return data, err
}

func (d *Deal) Json() ([]byte, error) {
	return json.Marshal(d)
}

func (d Deal) Find(restUrl string) (response BitrixRestResponseFind, err error) {

	if restUrl[len(restUrl)-1:] == "/" {
		restUrl = restUrl[:len(restUrl)-1]
	}

	url := fmt.Sprintf("%s/%s?filter[ORIGIN_ID]=%s", restUrl, findDeal, d.Id)

	if res, err := ExecReq("GET", url, nil); err == nil {
		return checkResponseFind(res)
	} else {
		return response, err
	}
}

func (d Deal) Add(restUrl string) (response BitrixRestResponseAdd, err error) {

	if restUrl[len(restUrl)-1:] == "/" {
		restUrl = restUrl[:len(restUrl)-1]
	}

	url := fmt.Sprintf("%s/%s", restUrl, addDeal)

	err = d.checkContact(restUrl)
	if err != nil {
		return
	}

	rd, err := d.getReader()
	if err != nil {
		return
	}

	if res, err := ExecReq("POST", url, rd); err == nil {
		return checkResponseAdd(res)
	} else {
		return response, err
	}
}

func (d Deal) Update(restUrl string, id string) (response BitrixRestResponseUpdate, err error) {
	if restUrl[len(restUrl)-1:] == "/" {
		restUrl = restUrl[:len(restUrl)-1]
	}

	url := fmt.Sprintf("%s/%s?id=%s", restUrl, updateDeal, id)

	err = d.checkContact(restUrl)
	if err != nil {
		return
	}

	rd, err := d.getReader()
	if err != nil {
		return
	}

	if res, err := ExecReq("POST", url, rd); err == nil {
		return checkResponseUpdate(res)
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
			id := response.Result
			d.ContactId = id
		} else {
			return err
		}
	} else {
		id := response.Result[0].ID
		d.ContactId = converter.String(id).Int()
	}

	return nil
}

func (d Deal) getReader() (rd io.Reader, err error) {

	data := make(map[string]map[string]string)

	deal := make(map[string]string)

	deal["ORIGIN_ID"] = d.Id
	deal["TITLE"] = d.Name

	if d.User != 0 {
		deal["ASSIGNED_BY_ID"] = fmt.Sprint(d.User)
	}

	deal["BEGINDATE"] = d.Date
	deal["CATEGORY_ID"] = fmt.Sprint(d.Category)
	deal["CONTACT_ID"] = fmt.Sprint(d.ContactId)

	for _, v := range d.UsersFields {
		deal[v.Id] = v.Value
	}

	for _, v := range d.UserFieldsPlurar {
		deal[v.Id] = fmt.Sprintf("[%s]", strings.Join(v.Value, ","))
	}

	data["fields"] = deal

	if content, err := json.Marshal(data); err == nil {
		return bytes.NewReader(content), err
	} else {
		return rd, err
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
	return d
}

func newDealFromOrder(order sql.Order) (Deal, error) {

	d := Deal{}

	offset := 2000

	d.Id = order.Id
	d.Name = order.Name
	d.User = converter.String(order.UserId).Int()
	d.Date = converter.SubtractionYearsOffset(order.Date, offset, "02.01.2006")

	isInternetOrder := converter.String(order.InternetOrder).BinaryTrue()

	d.Category = converter.GetCategoryInOrder(order.OrderType, isInternetOrder)

	d.ContactData = newContactFromClient(order.Client)

	if isInternetOrder {
		if stg, err := converter.GetInternetOrderStage(order.InternetOrderStage); err == nil {
			d.Stage = stg
		} else {
			return d, err
		}
	} else {
		d.Stage = converter.GetOrderStageByKind(order.OrderType)
	}

	d.Comment = order.Comment
	d.Sum = converter.String(order.Sum).Float32()

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594208144",
			Value: converter.GetOrderType(order.OrderType),
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594208880",
			Value: order.Dpp,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594209390",
			Value: order.DppOD,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594209319",
			Value: order.DppOS,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594209423",
			Value: order.Department,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594211273",
			Value: converter.SubtractionYearsOffset(order.ShipmentDate, 2000, "02.01.2006"),
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594273942",
			Value: order.Agreement,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594274024",
			Value: order.Stock,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1589971684",
			Value: order.Doctor,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594274325",
			Value: converter.String(order.SentSms).BoolAsString(),
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594274114",
			Value: converter.GetDeliveryType(order.DeliveryWay),
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594274276",
			Value: order.PickUpPoint,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594275086",
			Value: order.DeliveryAddress,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594275147",
			Value: converter.SubtractionYearsOffset(order.WantedDateShipment, 2000, "02.01.2006"),
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594275197",
			Value: order.DeliveryArea,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594275243",
			Value: converter.SubtractionYearsOffset(order.DeliveryTimeFrom, 2000, "02.01.2006T15:04:05"),
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594275305",
			Value: converter.SubtractionYearsOffset(order.DeliveryTimeTo, 2000, "02.01.2006T15:04:05"),
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594275360",
			Value: order.ExtraInfo,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594284459",
			Value: order.Prepaid,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594284797",
			Value: order.Prepayment,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594284830",
			Value: order.Credit,
		})

	segments := make([]string, 0, len(order.Segments))

	for _, v := range order.Segments {
		if sg, err := converter.GetBitrixSegment(v.Id); err == nil {
			segments = append(segments, sg)
		}

		for _, b := range v.Brands {
			if fieldName, ok := converter.GetNameBrandFieldNameForSegment(v.Id); ok {
				d.UsersFields = append(
					d.UsersFields,
					UserField{Id: fieldName,
						Value: b,
					})
			}
		}
	}

	d.UserFieldsPlurar = append(d.UserFieldsPlurar,
		UserFiledPlurarValue{
			Id:    "UF_CRM_1640090278334",
			Value: segments,
		})

	return d, nil
}

func newDealFromShipment(shipment sql.Shipment) (Deal, error) {

	d := Deal{}

	offset := 2000

	d.Id = shipment.Id
	d.Name = shipment.Name
	d.User = converter.String(shipment.User).Int()
	d.Date = converter.SubtractionYearsOffset(shipment.Date, offset, "02.01.2006")

	d.Category = 7

	d.ContactData = newContactFromClient(shipment.Client)

	d.Stage = "C7:WON"

	d.Comment = shipment.Comment
	d.Sum = converter.String(shipment.Sum).Float32()

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594208144",
			Value: converter.GetOrderType("1"),
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594209423",
			Value: shipment.Department,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594273942",
			Value: shipment.Agreement,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1594274024",
			Value: shipment.Stock,
		})

	d.UsersFields = append(
		d.UsersFields,
		UserField{Id: "UF_CRM_1589971684",
			Value: shipment.Doctor,
		})

	segments := make([]string, 0, len(shipment.Segments))

	for _, v := range shipment.Segments {
		if sg, err := converter.GetBitrixSegment(v.Id); err == nil {
			segments = append(segments, sg)
		}

		for _, b := range v.Brands {
			if fieldName, ok := converter.GetNameBrandFieldNameForSegment(v.Id); ok {
				d.UsersFields = append(
					d.UsersFields,
					UserField{Id: fieldName,
						Value: b,
					})
			}
		}
	}

	d.UserFieldsPlurar = append(d.UserFieldsPlurar,
		UserFiledPlurarValue{
			Id:    "UF_CRM_1640090278334",
			Value: segments,
		})

	return d, nil
}