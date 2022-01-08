package sql

import (
	"database/sql"
	"fmt"
	"strings"
)

type Execute func(params map[string]string, key string) (data []map[string]string, err error)

func ExecuteQuery(db *sql.DB, texts map[string]string) Execute {
	return func(params map[string]string, key string) (data []map[string]string, err error) {

		if text, ok := texts[key]; ok {
			for k, v := range params {
				text = strings.ReplaceAll(text, k, v)
			}
			return executeAndReadQuery(db, text)
		} else {
			err = fmt.Errorf("not found query '%s' in texts map", key)
			return
		}
	}
}

func executeAndReadQuery(db *sql.DB, text string) (data []map[string]string, err error) {

	if rows, err := db.Query(text); err == nil {
		return readRows(rows)
	} else {
		return data, err
	}
}

func readRows(rows *sql.Rows) (data []map[string]string, err error) {

	if cols, err := rows.Columns(); err == nil {
		return scanRowsToMap(rows, cols)
	} else {
		return data, err
	}
}

func scanRowsToMap(rows *sql.Rows, cols []string) (data []map[string]string, err error) {

	columns := createSliceForResult(cols)

	for rows.Next() {
		pretty := make(map[string]string)

		if err := rows.Scan(columns[:]...); err != nil {
			return data, err
		}

		for i := range columns {
			val := *columns[i].(*interface{})
			pretty[cols[i]] = getColumnValueAsString(val)
		}

		data = append(data, pretty)

	}

	return data, err
}

func createSliceForResult(cols []string) []interface{} {
	results := make([]interface{}, len(cols))
	for i := range results {
		results[i] = new(interface{})
	}
	return results
}

func getColumnValueAsString(val interface{}) (res string) {

	if val == nil {
		return "NULL"
	}
	switch v := val.(type) {
	case []byte:
		res = string(v)
	default:
		res = fmt.Sprintf("%v", v)
	}

	return res
}
