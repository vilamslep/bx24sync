package bitrix24

type Entity interface {
	Json() ([]byte, error)
}
