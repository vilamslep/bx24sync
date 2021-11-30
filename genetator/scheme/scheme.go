package scheme

import (
	"github.com/vi-la-muerto/bx24-service/genetator/enum"
)

const (
	Male enum.Season = iota
	Female
)

const (
	AgreeOnSMS enum.CallIneternetShop = iota
	AgreeOnCalls
	Disagree
)

const (
	IsClient enum.ContactType = iota
	IsSupplier
	IsOther
)

const (
	IsTrue  = "0x01"
	IsFalse = "0x00"
)

type Client struct {
	Name       string
	SecondName string
	LastName   string
	Gender     enum.Season
	Deleted    bool
}
