package sql

import (
	"io"
	"strings"
	"testing"
)

func Test_CreateScheme_CreateSchemeWithoutError(t *testing.T) {
	correctlyReader := correctlySchemeReader()

	if _, err := CreateScheme(correctlyReader); err != nil {
		t.Errorf("Fall.Creating scheme is falled. Error:%v\n", err)
	}
}

func Test_CreateScheme_ReaderError(t *testing.T) {
	if _, err := CreateScheme(ErrorReader{}); err == nil {
		t.Error("Must be return an error. Error expected is 'Reader error'\n")
	}
}

func Test_CreateScheme_JsonSyntaxError(t *testing.T) {
	rd := strings.NewReader("Not json")
	if _, err := CreateScheme(rd); err == nil {
		t.Error("Must be return an error. Error expected is 'json/SyntaxError'")
	}
}

func correctlySchemeReader() io.Reader {
	return strings.NewReader(`
	[
		{ 
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
