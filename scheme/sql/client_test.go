package sql

import (
	"fmt"
	"testing"
)

type ErrorReader struct{}

func (r ErrorReader) Read(b []byte) (n int, err error) {
	return 0, fmt.Errorf("Reader error")
}

func Test_convert_WithoutError(t *testing.T) {

	rightClient := Client{
		OriginId:                  "cc645ffc-9eb0-11e3-b1c9-00025553e653",
		Name:                      "Name",
		Birthday:                  "26.12.1969",
		Gender:                    "1",
		IsClient:                  "0x01",
		IsSuppler:                 "0x01",
		OtherRelation:             "0x00",
		IsRetireeOrDisabledPerson: "0x01",
		ConnectionWay:             "100",
		ThereIsContract:           "0x00",
		SendAds:                   "0x01",
		IsInternetClient:          "0x01",
		IsOfflineClient:           "0x01",
		IsClinicClient:            "0x00",
		DiscountClinicService:     "7",
		DiscountMedicalThings:     "0",
		DiscountRayban:            "10",
		Phone:                     "89091111111;",
		Email:                     "mail@mail.ru;",
	}

	item := map[string]string{
		"originId":                  "cc645ffc-9eb0-11e3-b1c9-00025553e653",
		"name":                      "Name",
		"birthday":                  "26.12.1969",
		"gender":                    "1",
		"isClient":                  "0x01",
		"isSuppler":                 "0x01",
		"otherRelation":             "0x00",
		"isRetireeOrDisabledPerson": "0x01",
		"connectionWay":             "100",
		"thereIsContract":           "0x00",
		"sendAds":                   "0x01",
		"isInternetClient":          "0x01",
		"isOfflineClient":           "0x01",
		"isClinicClient":            "0x00",
		"discountClinicService":     "7",
		"discountMedicalThings":     "0",
		"discountRayban":            "10",
		"phone":                     "89091111111;",
		"email":                     "mail@mail.ru;",
	}

	scheme := []Field{
		{Key: "originId", Required: true},
		{Key: "name", Required: true},
		{Key: "birthday", Required: true},
		{Key: "gender", Required: true},
		{Key: "isClient", Required: true},
		{Key: "isSuppler", Required: true},
		{Key: "otherRelation", Required: true},
		{Key: "isRetireeOrDisabledPerson", Required: true},
		{Key: "connectionWay", Required: true},
		{Key: "isRetireeOrDisabledPerson", Required: true},
		{Key: "thereIsContract", Required: true},
		{Key: "sendAds", Required: true},
		{Key: "isInternetClient", Required: true},
		{Key: "isOfflineClient", Required: true},
		{Key: "isClinicClient", Required: true},
		{Key: "discountClinicService", Required: true},
		{Key: "discountMedicalThings", Required: true},
		{Key: "discountRayban", Required: true},
		{Key: "phone", Required: true},
		{Key: "email", Required: true},
	}

	if client, err := convert(scheme, item); err != nil {
		t.Errorf("Converting error: %s ", err)
	} else if client != rightClient {
		t.Error("Clients must be same. But they are different")
	}
}

func Test_convert_NotFoundKeyError(t *testing.T) {
	data := make(map[string]string)
	data["id"] = "id"

	scheme := []Field{
		{Key: "originId", Required: true},
		{Key: "name", Required: true},
		{Key: "birthday", Required: true},
		{Key: "gender", Required: true},
		{Key: "isClient", Required: true},
		{Key: "isSuppler", Required: true},
		{Key: "otherRelation", Required: true},
		{Key: "isRetireeOrDisabledPerson", Required: true},
		{Key: "connectionWay", Required: true},
		{Key: "isRetireeOrDisabledPerson", Required: true},
		{Key: "thereIsContract", Required: true},
		{Key: "sendAds", Required: true},
		{Key: "isInternetClient", Required: true},
		{Key: "isOfflineClient", Required: true},
		{Key: "isClinicClient", Required: true},
		{Key: "discountClinicService", Required: true},
		{Key: "discountMedicalThings", Required: true},
		{Key: "discountRayban", Required: true},
		{Key: "phone", Required: true},
		{Key: "email", Required: true},
	}

	if _, err := convert(scheme, data); err == nil {
		t.Error("Must be return error. Error expected is  'not found item by key'")
	}
}
