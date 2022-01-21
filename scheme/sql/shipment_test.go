package sql

import "testing"

func Test_ConvertToShipment_WithoutError(t *testing.T) {
	t.Fail()
}

func getCorrectlyReception() Reception {
	return Reception{
		Id:         "1234",
		OriginId:   "uuid",
		Name:       "Reception",
		Date:       "28.06.4019",
		Department: "department",
		ClientId:   "clientId",
		Client: Client{
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
		},
		UserId: "111",
		AdditionnalFields: []AdditionalField{
			{Key: "Key1", Value: "Value1"},
			{Key: "Key2", Value: "Value2"},
			{Key: "Key3", Value: "Value3"},
			{Key: "Key4", Value: "Value4"},
			{Key: "Key5", Value: "Value5"},
			{Key: "Key6", Value: "Value6"},
		},
	}
}

func getCorrectltReceptionRaw() {

}
