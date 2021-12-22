package bitrix24

import "time"

type Deal struct {
	User int `json:"ASSIGNED_BY_ID"`
	Date time.Time `json:"BEGINDATE"`
	Category int `json:"CATEGORY_ID"`
	contact Contact
	ContactId int `json:"CONTACT_ID"`
	Sum float32 `json:"OPPORTUNITY"`
	Stage int `json:"STAGE_ID"`
	UsersFields []UserField
}

type UserField struct {
	Id string
	Value string
}