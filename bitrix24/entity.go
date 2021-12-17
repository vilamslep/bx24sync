package bitrix24

type Entity interface{
	LoadFromMap(map[string]string)
	Json() ([]byte, error) 
}