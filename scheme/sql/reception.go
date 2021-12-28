package sql

import "encoding/json"

type Reception struct {
	Id                string            `json:"id"`
	OriginId          string            `json:"originId"`
	Name              string            `json:"name"`
	Date              string            `json:"date"`
	Department        string            `json:"department"`
	ClientId          string            `json:"client"`
	Client            Client            `json:"clientData,omitempty"`
	UserId            string            `json:"userId"`
	AdditionnalFields []AdditionalField `json:"usersFields"`
}

func ConvertToReception(scheme []Field, data []map[string]string) ([]Reception, error) {
	res := make([]Reception, 0, len(data))

	for _, v := range data {
		if c, err := convertReception(scheme, v); err == nil {
			res = append(res, c)
		} else {
			return nil, err
		}
	}
	return res, nil
}

func convertReception(scheme []Field, data map[string]string) (c Reception, err error) {

	if err := checkByScheme(data, scheme); err != nil {
		return Reception{}, err
	}

	content, err := json.Marshal(&data)
	if err != nil {
		return Reception{}, err
	}

	err = json.Unmarshal(content, &c)

	return c, err
}
