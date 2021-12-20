package generator

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	bx24 "github.com/vi-la-muerto/bx24sync"
	schemes "github.com/vi-la-muerto/bx24sync/scheme/sql"
)

type executeDBQuery func(params map[string]string) (data []map[string]string, err error)

func HandlerWithDatabaseConnection(executeQuery executeDBQuery) bx24.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			return writeClientError(w, err, []byte("Don't manage to get body"))
		}

		id := "0x" + strings.ToUpper(string(body)[46:78])

		if data, err := executeQuery(map[string]string{"${client}": id}); err == nil {

			rd, err := os.OpenFile("clientScheme.json", os.O_RDONLY, 0666)

			if err != nil {
				return writeServerError(w, err)
			}

			scheme, err := schemes.CreateClientScheme(bufio.NewReader(rd))

			if err != nil {
				return writeServerError(w, err)
			}

			if content, err := schemes.ConvertThroughClient(data, scheme); err == nil {
				w.Write(content)
			} else {
				return writeServerError(w, err)
			}

		} else {
			return writeServerError(w, err)
		}
		return err
	}
}

func writeServerError(w http.ResponseWriter, err error) error {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	return err
}

func writeClientError(w http.ResponseWriter, err error, msg []byte) error {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(msg)
	return err
}

func ExecuteQuery(db *sql.DB, text string) executeDBQuery {
	return func(params map[string]string) (data []map[string]string, err error) {
		for k, v := range params {
			text = strings.ReplaceAll(text, k, v)
		}
		return executeAndReadQuery(db, text)
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

	pretty := make(map[string]string)
	for rows.Next() {

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
