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
	Phone                 []ContactData `json:"PHONE"`
	Email                 []ContactData `json:"EMAIL"`
}

type ContactData struct {
	ValueType string `json:"VALUE_TYPE"`
	Value     string `json:"VALUE"`
	TypeId    string `json:"TYPE_ID"`
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

	phones := data["Phone"]

	sPhone := strings.Split(phones, ";")

	for _, phone := range sPhone {
		if phone != "" {
			contData := ContactData{
				ValueType: "WORK",
				Value:     phone,
				TypeId:    "PHONE",
			}

			s.Phone = append(s.Phone, contData)
		}
	}

	emails := data["Email"]

	sEmail := strings.Split(emails, ";")

	for _, email := range sEmail {
		if email != "" {
			contData := ContactData{
				ValueType: "WORK",
				Value:     email,
				TypeId:    "EMAIL",
			}

			s.Email = append(s.Phone, contData)
		}
	}
}


type GeneratorConfig struct {
	DB        DataBaseConfig `json:"DB"`
	Web       WebConfig      `json:"Web"`
	QueryPath string         `json:"QueryPath"`
}

type DataBaseConfig struct {
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	User     string `json:"User"`
	Password string `json:"Password"`
	Database string `json:"Database"`
}

type WebConfig struct {
	Host string `json:"Host"`
	Port int    `json:"Port"`
}