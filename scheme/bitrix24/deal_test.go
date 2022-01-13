package bitrix24

import (
	"encoding/json"
	"testing"
)

func Test_AddShipment(t *testing.T) {
	data := getShipmentRaw()

	entity, err := NewDealFromJson([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	restUrl := "https://portal.optic-center.ru/rest/475/vq2fo1npkvh2m1fd"
	response, err := entity.Find(restUrl)

	if err != nil {
		t.Fatal(err)
	}

	if response.Total == 0 {
		if _, err := entity.Add(restUrl); err != nil {
			t.Fatal(err)
		}
	} else {
		id := response.Result[0].ID

		if r, err := entity.Update(restUrl, id); err != nil {
			t.Fatal(err)
		} else {
			if !r.Result {
				t.Fatal("can't to update entity")
			}
		}
	}
}

func Test_UnmarshalResponseAdd(t *testing.T) {

	res := getDataResponseData()

	response := BitrixRestResponseAdd{}

	if err := json.Unmarshal([]byte(res), &response); err != nil {
		t.Fatal(err)
	} else {
		t.Log(response)
	}

}

func getDataResponseData() string {
	return `{
		"result":1776521,
		"time":{
			"start":1642060235.149831,
			"finish":1642060235.976128,
			"duration":0.8262970447540283,
			"processing":0.7797980308532715,
			"date_start":"2022-01-13T10:50:35+03:00",
			"date_finish":"2022-01-13T10:50:35+03:00"
			}
			}`
}

func getShipmentRaw() string {
	return ``
}
