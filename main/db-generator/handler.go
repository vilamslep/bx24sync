package main

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	bx24 "github.com/vi-la-muerto/bx24sync"
	schemeConf "github.com/vi-la-muerto/bx24sync/scheme/sql"
	"github.com/vi-la-muerto/bx24sync/sql"
)

func HandlerClientWithDatabaseConnection(executeQuery sql.Execute) bx24.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			return writeBadRequests(w, err, []byte("Don't manage to get body"))
		}

		id := "0x" + strings.ToUpper(string(body)[46:78])

		if data, err := executeQuery(map[string]string{"${client}": id}, "client"); err == nil {

			rd, err := os.OpenFile("client.json", os.O_RDONLY, 0666)

			if err != nil {
				return writeServerError(w, err)
			}

			scheme, err := schemeConf.CreateScheme(bufio.NewReader(rd))

			if err != nil {
				return writeServerError(w, err)
			}

			if clients, err := schemeConf.ConvertToClients(scheme, data); err == nil {
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

func HandlerReceptionWithDatabaseConnection(executeQuery sql.Execute) bx24.HandlerFunc {
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

			scheme, err := schemeConf.CreateScheme(bufio.NewReader(rd))

			if err != nil {
				return writeServerError(w, err)
			}

			if repections, err := schemeConf.ConvertToReception(scheme, data); err == nil {

				repections, err = setClientReception(repections, executeQuery)

				if err != nil {
					return writeServerError(w, err)
				}

				repections, err = setReceptionPropertyes(repections, executeQuery)

				if err != nil {
					return writeServerError(w, err)
				}

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

func HandlerOrderWithDatabaseConnection(executeQuery sql.Execute) bx24.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			return writeBadRequests(w, err, []byte("Don't manage to get body"))
		}

		id := "0x" + strings.ToUpper(string(body)[46:78])

		args := map[string]string{
			"${order}": id,
		}

		if data, err := executeQuery(args, "order"); err == nil {

			rd, err := os.OpenFile("order.json", os.O_RDONLY, 0666)

			if err != nil {
				return writeServerError(w, err)
			}

			scheme, err := schemeConf.CreateScheme(bufio.NewReader(rd))

			if err != nil {
				return writeServerError(w, err)
			}

			if orders, err := schemeConf.ConvertToOrders(scheme, data); err == nil {

				orders, err = setClientOrders(orders, executeQuery)

				if err != nil {
					return writeServerError(w, err)
				}

				orders, err = setOrderSegments(orders, executeQuery)

				if err != nil {
					return writeServerError(w, err)
				}

				content, err := json.Marshal(orders)
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

func HandlerShipmentWithDatabaseConnection(executeQuery sql.Execute) bx24.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			return writeBadRequests(w, err, []byte("Don't manage to get body"))
		}

		id := "0x" + strings.ToUpper(string(body)[46:78])

		args := map[string]string{
			"${shipment}": id,
		}

		if data, err := executeQuery(args, "shipment"); err == nil {

			rd, err := os.OpenFile("shipment.json", os.O_RDONLY, 0666)

			if err != nil {
				return writeServerError(w, err)
			}

			scheme, err := schemeConf.CreateScheme(bufio.NewReader(rd))

			if err != nil {
				return writeServerError(w, err)
			}

			if shipment, err := schemeConf.ConvertToShipment(scheme, data); err == nil {

				shipment, err = setClientShipment(shipment, executeQuery)

				if err != nil {
					return writeServerError(w, err)
				}

				shipment, err = setShipmentSegments(shipment, executeQuery)

				if err != nil {
					return writeServerError(w, err)
				}

				content, err := json.Marshal(shipment)
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
