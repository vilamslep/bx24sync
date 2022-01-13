package bitrix24

type Entity interface {
	Add(restUrl string) (BitrixRestResponseAdd, error)
	Find(restUrl string) (BitrixRestResponseFind, error)
	Update(restUrl string, id string) (BitrixRestResponseUpdate, error)
}