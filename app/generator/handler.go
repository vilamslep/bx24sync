package generator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	bx24 "github.com/vi-la-muerto/bx24sync"
	schemes "github.com/vi-la-muerto/bx24sync/scheme/sql"
)

func HandlerClientWithDatabaseConnection(executeQuery bx24.Execute) bx24.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			return writeBadRequests(w, err, []byte("Don't manage to get body"))
		}

		id := "0x" + strings.ToUpper(string(body)[46:78])

		if data, err := executeQuery(map[string]string{"${client}": id}, "client"); err == nil {

			rd, err := os.OpenFile("client.json", os.O_RDONLY, 0666)

			if err != nil {
				return writeServerError(w, err)
			}

			scheme, err := schemes.CreateScheme(bufio.NewReader(rd))

			if err != nil {
				return writeServerError(w, err)
			}

			if clients, err := schemes.ConvertToClients(scheme, data); err == nil {
				content, err := json.Marshal(clients)
				if err != nil {
					return writeServerError(w, err)
				}

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

func HandlerReceptionWithDatabaseConnection(executeQuery bx24.Execute) bx24.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			return writeBadRequests(w, err, []byte("Don't manage to get body"))
		}

		id := "0x" + strings.ToUpper(string(body)[46:78])

		args := map[string]string{
			"${reception}": id,
		}

		if data, err := executeQuery(args, "reception"); err == nil {

			rd, err := os.OpenFile("reception.json", os.O_RDONLY, 0666)

			if err != nil {
				return writeServerError(w, err)
			}

			scheme, err := schemes.CreateScheme(bufio.NewReader(rd))

			if err != nil {
				return writeServerError(w, err)
			}

			if repections, err := schemes.ConvertToReception(scheme, data); err == nil {

				repections, err = setClient(repections, executeQuery)

				if err != nil { return writeServerError(w, err) }

				repections, err = setPropertyes(repections, executeQuery)

				if err != nil { return writeServerError(w, err) }

				content, err := json.Marshal(repections)
				if err != nil {
					return writeServerError(w, err)
				}
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

func HandlerOrderWithDatabaseConnection(executeQuery bx24.Execute) bx24.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()

		writeResponse(w, http.StatusForbidden, nil)
		return nil
	}
}

func HandlerShipmentWithDatabaseConnection(executeQuery bx24.Execute) bx24.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()

		writeResponse(w, http.StatusForbidden, nil)
		return nil
	}
}

func writeServerError(w http.ResponseWriter, err error) error {
	writeResponse(w, http.StatusInternalServerError, []byte(err.Error()))
	return err
}

func writeBadRequests(w http.ResponseWriter, err error, msg []byte) error {
	writeResponse(w, http.StatusBadRequest, msg)
	return err
}

func writeResponse(w http.ResponseWriter, status int, msg []byte) {
	w.WriteHeader(status)
	w.Write(msg)
}

func setClient(res []schemes.Reception, executeQuery bx24.Execute) ([]schemes.Reception, error) {
	result := make([]schemes.Reception, 0, len(res))

	rd, err := os.OpenFile("client.json", os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	scheme, err := schemes.CreateScheme(rd)
	if err != nil {
		return nil, err
	}

	for _, reception := range res {

		id := fmt.Sprintf("0x%s", reception.ClientId)

		args := map[string]string{
			"${client}": id,
		}

		data, err := executeQuery(args, "client")

		if res, err := schemes.ConvertToClients(scheme, data); err == nil {
			if len(res) > 0 {
				reception.Client = res[0]
			}
		} else {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		result = append(result, reception)
	}
	return result, err
}

func setPropertyes(res []schemes.Reception, executeQuery bx24.Execute) ([]schemes.Reception, error) {

	result := make([]schemes.Reception, 0, len(res))

	rd, err := os.OpenFile("reception_propertyes.json", os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	scheme, err := schemes.CreateScheme(rd)
	if err != nil {
		return nil, err
	}

	for _, reception := range res {

		args := map[string]string{
			"${reception}": fmt.Sprintf("0x%s", reception.Id),
		}
		data, err := executeQuery(args, "reception_propertyes")

		if res, err := schemes.ConvertToAdditionalFields(scheme, data); err == nil {
			reception.AdditionnalFields = res
		} else {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		result = append(result, reception)
	}
	return result, err
}
