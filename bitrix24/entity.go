package bitrix24

import (
	"strconv"
	"strings"
)

type Gender int
type ConnectionWay int

const (
	Male   Gender = 263
	Female Gender = 264
)

const (
	SMS   ConnectionWay = 515
	Phone ConnectionWay = 516
)

const (
	binaryTrue = "\x01"
)

type ContactData struct {
	ValueType string `json:"VALUE_TYPE"`
	Value     string `json:"VALUE"`
	TypeId    string `json:"TYPE_ID"`
}

type Contact struct {
	Id                    string        `json:"ORIGIN_ID"`
	Name                  string        `json:"NAME"`
	SecondName            string        `json:"SECOND_NAME"`
	LastName              string        `json:"LAST_NAME"`
	Birthday              string        `json:"BIRTHDATE"`
	Gender                Gender        `json:"UF_CRM_1587720922"`
	ConnectionWay         ConnectionWay `json:"UF_CRM_1612972723496"`
	ThereIsContract       bool          `json:"UF_CRM_1584964931121"`
	SendAds               bool          `json:"UF_CRM_1584970848775"`
	NotSendAds            bool          `json:"UF_CRM_1636547489168"`
	IsInternetClient      bool          `json:"UF_CRM_1585077731269"`
	IsOfflineClient       bool          `json:"UF_CRM_1585077757363"`
	IsClinicClient        bool          `json:"UF_CRM_1585077767559"`
	DiscountClinicService uint8         `json:"UF_CRM_1587720973"`
	DiscountMedicalThings uint8         `json:"UF_CRM_1587720952"`
	DiscountRayban        uint8         `json:"UF_CRM_1587720989"`
	Phone                 []ContactData `json:"PHONE"`
	Email                 []ContactData `json:"EMAIL"`
}

func (s *Contact) TransoftFromMap(data map[string]string) {

	s.Id = data["originId"]
	s.Birthday = data["birthday"]

	s.transforName(data["name"])

	s.DiscountMedicalThings = stringToUint8(data["discountMedicalThings"])
	s.DiscountRayban = stringToUint8(data["discountRayban"])
	s.DiscountClinicService = stringToUint8(data["discountClinicService"])

	s.Gender = stringToGenderCode(data["gender"])
	s.ConnectionWay = stringToConnectionWay(data["connectionWay"])

	s.ThereIsContract = isBinaryTrue(data["thereIsContract"])
	s.IsOfflineClient = isBinaryTrue(data["isOfflineClient"])
	s.IsInternetClient = isBinaryTrue(data["isInternetClient"])
	s.SendAds = isBinaryTrue(data["sendAdds"])
	s.NotSendAds = !s.SendAds
	s.IsClinicClient = isBinaryTrue(data["isClinicClient"])

	s.Phone = splitStringToContactDateSlice(data["Phone"], ";", "PHONE", "PHONE")

	s.Email = splitStringToContactDateSlice(data["Email"], ";", "EMAIL", "EMAIL")
}

func (s *Contact) transforName(val string) {

	valWithoutNum := val
	for _, v := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"} {
		valWithoutNum = strings.ReplaceAll(valWithoutNum, string(v), "")
	}

	pathName := strings.Split(valWithoutNum, " ")

	lengh := len(pathName)
	if lengh >= 3 {
		s.LastName, s.Name, s.SecondName = pathName[0], pathName[1], pathName[2]
	} else if lengh == 2 {
		s.LastName, s.Name = pathName[0], pathName[1]
	} else {
		s.Name = val
	}
}

func isBinaryTrue(value string) bool {
	return value == binaryTrue
}

func stringToUint8(val string) uint8 {

	result := uint8(0)

	if value, err := strconv.Atoi(val); err == nil {
		result = uint8(value)
	}
	return result
}

func stringToGenderCode(val string) Gender {

	result := Gender(0)

	switch Gender(stringToUint8(val)) {
	case 0:
		result = Male
	case 1:
		result = Female
	}

	return result
}

func stringToConnectionWay(val string) ConnectionWay {

	result := ConnectionWay(0)
	switch ConnectionWay(stringToUint8(val)) {
	case 0:
		result = SMS
	case 1:
		result = Phone
	}

	return result
}

func splitStringToContactDateSlice(val string, sep string, valType string, typeId string) []ContactData {

	if len(val) == 0 {
		return nil
	}

	slice := strings.Split(val, sep)

	target := make([]ContactData, 0)

	for _, piece := range slice {
		if piece == "" {
			continue
		}

		item := ContactData{
			ValueType: valType,
			Value:     piece,
			TypeId:    typeId,
		}

		target = append(target, item)

	}

	return target
}
