package scheme

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
	byteTrue = "\x01"
)

type Client struct {
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
	Phone                 string        `json:"PHONE"`
	Email                 string        `json:"EMAIL"`
}

func (s *Client) TransoftFromMap(data map[string]string) {
	s.Id = data["originId"]

	if value, err := strconv.Atoi(data["gender"]); err == nil {
		switch Gender(value) {
		case 0:
			s.Gender = Male
		case 1:
			s.Gender = Female
		}
	}

	if value, err := strconv.Atoi(data["discountMedicalThings"]); err == nil {
		s.DiscountMedicalThings = uint8(value)
	}

	if value, err := strconv.Atoi(data["discountRayban"]); err == nil {
		s.DiscountRayban = uint8(value)
	}

	if value, err := strconv.Atoi(data["discountClinicService"]); err == nil {
		s.DiscountClinicService = uint8(value)
	}

	if value, err := strconv.Atoi(data["connectionWay"]); err == nil {
		switch ConnectionWay(value) {
		case 0:
			s.ConnectionWay = SMS
		case 1:
			s.ConnectionWay = Phone
		}
	}

	s.ThereIsContract = byteTrue == data["thereIsContract"]
	s.IsOfflineClient = byteTrue == data["isOfflineClient"]
	s.IsInternetClient = byteTrue == data["isInternetClient"]
	s.SendAds = byteTrue == data["sendAdds"]
	s.NotSendAds = !s.SendAds
	s.IsClinicClient = byteTrue == data["isClinicClient"]

	pathName := strings.Split(data["name"], " ")

	if len(pathName) == 3 {
		s.LastName = pathName[0]
		s.Name = pathName[1]
		s.SecondName = pathName[2]
	} else if len(pathName) == 2 {
		s.LastName = pathName[0]
		s.Name = pathName[1]
	} else {
		s.Name = data["name"]
	}

	s.Birthday = data["birthday"]
}
