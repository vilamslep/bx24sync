package sql

type AdditionalField struct {
	Key   string
	Value string
}

func ConvertToAdditionalFields(scheme []Field, data []map[string]string) ([]AdditionalField, error) {
	res := make([]AdditionalField, 0, len(data))

	for _, v := range data {
		if c, err := convertAdditionalField(v, scheme); err == nil {
			res = append(res, c)
		} else {
			return nil, err
		}
	}
	return res, nil
}

func convertAdditionalField(data map[string]string, scheme []Field) (c AdditionalField, err error) {

	if err := checkByScheme(data, scheme); err != nil {
		return AdditionalField{}, err
	}

	c.Key = data["id"]
	switch data["valType"]{
	case "2": //bool
		c.Value = data["valBool"]
	case "3": //float
		c.Value = data["valDigit"]
	case "4": //time
		c.Value = data["valTime"]
	case "5": //str
		c.Value = data["valStr"]
	case "8": //ref 
		c.Value = data["valRef"]
	}

	return c, err
}

type Segment struct {
	Type string 
	Brand string
}