package sql

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

type ErrorReader struct{}

func (r ErrorReader) Read(b []byte) (n int, err error) {
	return 0, fmt.Errorf("Reader error")
}

func Test_CreateClientScheme_Success(t *testing.T) {
	if _, err := CreateScheme(getReaderSchemeForCorrectlyScheme()); err != nil {
		t.Error(err)
	}
}
func Test_CreateClientScheme_Error_ReaderError(t *testing.T) {
	if _, err := CreateScheme(ErrorReader{}); err == nil {
		t.Error("Must be return error. Error expected is 'Reader error'")
	}
}
func Test_CreateClientScheme_Error_UnmarshalError(t *testing.T) {
	if _, err := CreateScheme(getReaderSchemeNotJson()); err == nil {
		t.Error("Must be return error. Error expected is 'json/SyntaxError'")
	}
}

func Test_CheckByScheme_Success(t *testing.T) {

	testData, err := getCorrectlyTestDataFromFile()

	if err != nil {
		t.Fatal(err)
	}

	scheme, err := CreateScheme(getReaderSchemeForCorrectlyScheme())

	if err != nil {
		t.Error(err)
	}

	for _, v := range testData {
		if err := checkByScheme(v, scheme); err != nil {
			t.Error(err)
		}
	}
}
func Test_CheckByScheme_Error_NotFoundKey(t *testing.T) {

	testData := getUncorrectlyTestData()

	scheme, err := CreateScheme(getReaderSchemeForCorrectlyScheme())

	if err != nil {
		t.Error(err)
	}

	for _, v := range testData {
		if err := checkByScheme(v, scheme); err == nil {
			t.Error("Must be return error. Error expected is  'not found item by key'")
		}
	}
}

func Test_ConvertToClient_Success(t *testing.T) {
	rightClient := getRightClient()

	item := getCorrectlyTestData()

	scheme, err := CreateScheme(getReaderSchemeForCorrectlyScheme())

	if err != nil {
		t.Errorf("Error with correctly scheme reader: %s", err.Error())
	}

	if client, err := convert(scheme, item); err != nil {
		t.Errorf("Converting error: %s ", err)
	} else if client != rightClient {
		t.Error("Clients must be same. But they are different")
	}
}
func Test_ConvertToClient_Error_KeyNotFound(t *testing.T) {

	data := make(map[string]string)
	data["id"] = "id"

	scheme, err := CreateScheme(getReaderSchemeForCorrectlyScheme())

	if err != nil {
		t.Errorf("Error with correctly scheme reader: %s", err.Error())
	}

	if _, err := convert(scheme, data); err == nil {
		t.Error("Must be return error. Error expected is  'not found item by key'")
	}
}

func Test_ConvertThroughClient_Success(t *testing.T) {

	items, err := getCorrectlyTestDataFromFile()

	if err != nil {
		t.Errorf("Getting test data falied: %s", err.Error())
	}
	scheme, err := CreateScheme(getReaderSchemeForCorrectlyScheme())

	if err != nil {
		t.Errorf("Error with correctly scheme reader: %s", err.Error())
	}

	if _, err := ConvertToClients(scheme, items); err != nil {
		t.Errorf("Converting error: %s ", err)
	}
}

//test data
func getCorrectlyTestDataFromFile() ([]map[string]string, error) {

	testData := make([]map[string]string, 0)

	data := make([]string, 0)
	content, err := os.ReadFile("test_data.json")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &data)

	if err != nil {
		return nil, err
	}

	for _, v := range data {
		item := make([]map[string]string, 0)
		err = json.Unmarshal([]byte(v), &item)

		if err != nil {
			return nil, err
		}
		testData = append(testData, item...)
	}

	return testData, nil
}

func getCorrectlyTestData() map[string]string {

	item := make(map[string]string)
	item["originId"] = "cc645ffc-9eb0-11e3-b1c9-00025553e653"
	item["name"] = "Name"
	item["birthday"] = "26.12.1969"
	item["gender"] = "1"
	item["isClient"] = "0x01"
	item["isSuppler"] = "0x01"
	item["otherRelation"] = "0x00"
	item["isRetireeOrDisabledPerson"] = "0x01"
	item["connectionWay"] = "100"
	item["thereIsContract"] = "0x00"
	item["sendAds"] = "0x01"
	item["isInternetClient"] = "0x01"
	item["isOfflineClient"] = "0x01"
	item["isClinicClient"] = "0x00"
	item["discountClinicService"] = "7"
	item["discountMedicalThings"] = "0"
	item["discountRayban"] = "10"
	item["phone"] = "89090884922;"
	item["email"] = ""

	return item
}

func getRightClient() (c Client) {
	return Client{
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
		Phone:                     "89090884922;",
		Email:                     "",
	}
}

func getReaderSchemeForCorrectlyScheme() io.Reader {
	return strings.NewReader(`[{
        "Key": "originId",
        "Required": true
    },
    {
        "Key": "name",
        "Required": true
    },
    {
        "Key": "birthday",
        "Required": true
    },
    {
        "Key": "gender",
        "Required": true
    },
    {
        "Key": "isClient",
        "Required": true
    },
    {
        "Key": "isSuppler",
        "Required": true
    },
    {
        "Key": "otherRelation",
        "Required": true
    },
    {
        "Key": "isRetireeOrDisabledPerson",
        "Required": true
    },
    {
        "Key": "connectionWay",
        "Required": true
    },
    {
        "Key": "isRetireeOrDisabledPerson",
        "Required": true
    },
    {
        "Key": "thereIsContract",
        "Required": true
    },
    {
        "Key": "sendAds",
        "Required": true
    },
    {
        "Key": "isInternetClient",
        "Required": true
    },
    {
        "Key": "isOfflineClient",
        "Required": true
    },
    {
        "Key": "isClinicClient",
        "Required": true
    },
    {
        "Key": "discountClinicService",
        "Required": true
    },
    {
        "Key": "discountMedicalThings",
        "Required": true
    },
    {
        "Key": "discountRayban",
        "Required": true
    },
    {
        "Key": "phone",
        "Required": true
    },
    {
        "Key": "email",
        "Required": true
    }
]`)
}

func getReaderSchemeNotJson() io.Reader {
	return strings.NewReader("Not json")
}

//not defined some fields
func getUncorrectlyTestData() (res []map[string]string) {
	item := make(map[string]string)
	item["originId"] = "cc645ffc-9eb0-11e3-b1c9-00025553e653"
	item["gender"] = "1"
	item["isClient"] = "0x01"
	item["otherRelation"] = "0x00"
	item["connectionWay"] = "100"
	item["thereIsContract"] = "0x00"
	item["isOfflineClient"] = "0x01"
	item["isClinicClient"] = "0x00"
	item["discountClinicService"] = "7"
	item["discountRayban"] = "0"
	item["phone"] = "89090884922;"
	item["email"] = ""

	return append(res, item)
}
