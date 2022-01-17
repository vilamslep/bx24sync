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

func Test_Client(t *testing.T) {
	testID := 1
	fmt.Println("Test creating a scheme from raw data")
	{
		fmt.Printf("Test ID: %d. Test correctly way.\n", testID)
		{
			correctlyReader := getReaderSchemeForCorrectlyScheme()
			if _, err := CreateScheme(correctlyReader); err != nil {
				t.Errorf("Fall.Creating scheme is falled. Error:%v\n", err)
			} else {
				fmt.Printf("Test ID: %d. SUCCESS \n", testID)
			}
		}

		testID++
		fmt.Printf("Test ID: %d. Test error cases.\n", testID)
		{
			if _, err := CreateScheme(ErrorReader{}); err == nil {
				t.Error("Must be return an error. Error expected is 'Reader error'\n")
			} else if _, err := CreateScheme(getReaderSchemeNotJson()); err == nil {
				t.Error("Must be return an error. Error expected is 'json/SyntaxError'")
			}
		}
	}

	fmt.Println("Test checking a scheme from raw data")
	{
		testID++
		fmt.Printf("Test ID: %d. Test correctly way.\n", testID)
		{
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

		testID++
		fmt.Printf("Test ID: %d. Test error cases\n", testID)
		{
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
	}

	fmt.Println("Test convertion raw to client")
	{
		testID++
		fmt.Printf("Test ID %d. Test correctly way\n", testID)
		{
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

		testID++
		fmt.Printf("Test ID %d. Test error cases\n", testID)
		{
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
		testID++
		fmt.Printf("Test ID %d. Test public method for converting", testID)
		{
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
	}
}

// test data
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

	return map[string]string{
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
		Phone:                     "89091111111;",
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
	item["phone"] = "89091111111;"
	item["email"] = "mail@mail.ru"

	return append(res, item)
}
