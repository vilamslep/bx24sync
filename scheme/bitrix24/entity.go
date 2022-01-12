package bitrix24

type Entity interface {
	Add(restUrl string) (BitrixRestResponse, error)
	Find(restUrl string) (BitrixRestResponse, error)
	Update(restUrl string, id string) (BitrixRestResponseUpdate, error)
}