package bitrix24

type Entity interface {
	Add() (BitrixRestResponse, error)
	// Get() (BitrixRestResponse, error)
	Find() (BitrixRestResponse, error)
	Update() (BitrixRestResponse, error)
}