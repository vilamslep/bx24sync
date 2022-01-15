package bitrix24

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type gettingData func(io.Reader) ([][]byte, error)

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

func getShipmentRaw() string {
	return ``
}

func Test_InternetOrder(t *testing.T) {
	id := getInternetOrder()

	creating := GetDealFromRawAsOrder
	rd := strings.NewReader(id)
	url := "http://95.78.157.195:25473/order"
	if response, err := createAndExecRequest("POST", url, rd); err == nil {
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {

			content, err := io.ReadAll(response.Body)
			if err != nil {
				t.Errorf("reading response.%s", err.Error())
			}
			t.Errorf("bad response from generator. Reponse: %s", string(content))
			return
		}

		data, err := convertDataForCrm(rd, creating)

		if err != nil {
			t.Error(err)
		}

		for _, v := range data {
			t.Log(string(v))
		}

	} else {
		t.Error(err)
	}
}

func getInternetOrder() string {
	return `{"#",dbdba9ef-a[34/1985]8330-4aaf831b4c73,367:b3c0a4bf015829f711ec609d1950adff}`
}

func createAndExecRequest(method string, url string, rd io.Reader) (*http.Response, error) {

	if req, err := http.NewRequest(method, url, rd); err == nil {
		client := http.Client{Timeout: time.Second * 300}
		return client.Do(req)
	} else {
		return nil, err
	}
}

func convertDataForCrm(r io.Reader, creating gettingData) (data [][]byte, err error) {
	return creating(r)
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
