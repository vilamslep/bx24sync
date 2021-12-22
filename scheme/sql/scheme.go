package sql

import (
	"encoding/json"
	"fmt"
	"io"
)

type Field struct {
	Key      string
	Required bool
}

func CreateScheme(reader io.Reader) (scheme []Field, err error) {
	if content, err := io.ReadAll(reader); err == nil {
		if err := json.Unmarshal(content, &scheme); err != nil {
			return scheme, err
		}
	} else {
		return scheme, err
	}
	return scheme, err
}

func checkByScheme(data map[string]string, scheme []Field) (err error) {
	for _, f := range scheme {
		if _, ok := data[f.Key]; !ok {
			if f.Required {
				err = fmt.Errorf("not found item by key '%s'", f.Key)
				break
			}
		}
	}
	return err
}