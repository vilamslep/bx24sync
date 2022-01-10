package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	bx24 "github.com/vi-la-muerto/bx24sync"
	schemeConf "github.com/vi-la-muerto/bx24sync/scheme/sql"
	sqlI "github.com/vi-la-muerto/bx24sync/sql"
	mssql "github.com/vi-la-muerto/bx24sync/sql/mssql"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("execution error %v", err)
	}
}

func run() (err error) {

	config := bx24.NewGeneratorConfigFromEnv()

	router := bx24.NewRouter(os.Stdout, os.Stderr, true)

	db, err := mssql.GetDatabaseConnection(config.DB)

	if err != nil {
		return err
	}
	if err := settingRouter(router, config.CheckInput, db, config.StorageQueryTxt); err != nil {
		log.Fatalf("can't to setting router. Err: %+v", err)
	}

	server := &http.Server{
		Addr:    config.Web.String(),
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Can't start to listener: %s\n", err)
	}

	return err
}

func settingRouter(r bx24.Router, enableCheckInput bool, db *sql.DB, queryTextsDir string) (err error) {

	var checkInputFunc bx24.CheckInput = nil

	if enableCheckInput {
		checkInputFunc = bx24.DefaultCheckInput
	}

	allowsMethods := []string{"POST"}

	txtCl, err := getQueryText(fmt.Sprintf("%s/client.sql", queryTextsDir))
	if err != nil {
		return
	}

	txtRecept, err := getQueryText(fmt.Sprintf("%s/reception.sql", queryTextsDir))
	if err != nil {
		return
	}

	txtReceptProp, err := getQueryText(fmt.Sprintf("%s/reception_propertyes.sql", queryTextsDir))
	if err != nil {
		return
	}

	txtOrder, err := getQueryText(fmt.Sprintf("%s/order.sql", queryTextsDir))
	if err != nil {
		return
	}

	txtOrderSegment, err := getQueryText(fmt.Sprintf("%s/order_segments.sql", queryTextsDir))
	if err != nil {
		return
	}

	txtShipment, err := getQueryText(fmt.Sprintf("%s/shipment.sql", queryTextsDir))
	if err != nil {
		return
	}

	txtShipmentSegment, err := getQueryText(fmt.Sprintf("%s/shipment_segments.sql", queryTextsDir))
	if err != nil {
		return
	}
	textsQuery := map[string]string{
		"client":               string(txtCl),
		"reception":            string(txtRecept),
		"reception_propertyes": string(txtReceptProp),
		"order":                string(txtOrder),
		"order_segments":       string(txtOrderSegment),
		"shipment":             string(txtShipment),
		"shipment_segments":    string(txtShipmentSegment),
	}

	r.AddMethod(
		bx24.HttpMethod{
			Path:         "/client",
			Handler:      HandlerClientWithDatabaseConnection(sqlI.ExecuteQuery(db, textsQuery)),
			CheckInput:   checkInputFunc,
			AllowMethods: allowsMethods,
		})

	r.AddMethod(
		bx24.HttpMethod{
			Path:         "/reception",
			Handler:      HandlerReceptionWithDatabaseConnection(sqlI.ExecuteQuery(db, textsQuery)),
			CheckInput:   checkInputFunc,
			AllowMethods: allowsMethods,
		})

	r.AddMethod(
		bx24.HttpMethod{
			Path:         "/order",
			Handler:      HandlerOrderWithDatabaseConnection(sqlI.ExecuteQuery(db, textsQuery)),
			CheckInput:   checkInputFunc,
			AllowMethods: allowsMethods,
		})

	r.AddMethod(
		bx24.HttpMethod{
			Path:         "/shipment",
			Handler:      HandlerShipmentWithDatabaseConnection(sqlI.ExecuteQuery(db, textsQuery)),
			CheckInput:   checkInputFunc,
			AllowMethods: allowsMethods,
		})

	return err
}

func getQueryText(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func HandlerClientWithDatabaseConnection(executeQuery sqlI.Execute) bx24.HandlerFunc {
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

func HandlerReceptionWithDatabaseConnection(executeQuery sqlI.Execute) bx24.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		body, err := io.ReadAll(r.Body)

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

func HandlerOrderWithDatabaseConnection(executeQuery sqlI.Execute) bx24.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		body, err := io.ReadAll(r.Body)

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

func HandlerShipmentWithDatabaseConnection(executeQuery sqlI.Execute) bx24.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		body, err := io.ReadAll(r.Body)

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

func setClientReception(res []schemeConf.Reception, executeQuery sqlI.Execute) ([]schemeConf.Reception, error) {
	result := make([]schemeConf.Reception, 0, len(res))

	rd, err := os.OpenFile("client.json", os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	scheme, err := schemeConf.CreateScheme(rd)
	if err != nil {
		return nil, err
	}

	for _, reception := range res {

		id := fmt.Sprintf("0x%s", reception.ClientId)

		args := map[string]string{
			"${client}": id,
		}

		data, err := executeQuery(args, "client")

		if res, err := schemeConf.ConvertToClients(scheme, data); err == nil {
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

func setClientOrders(res []schemeConf.Order, executeQuery sqlI.Execute) ([]schemeConf.Order, error) {
	result := make([]schemeConf.Order, 0, len(res))

	rd, err := os.OpenFile("client.json", os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	scheme, err := schemeConf.CreateScheme(rd)
	if err != nil {
		return nil, err
	}

	for _, order := range res {

		id := fmt.Sprintf("0x%s", order.ClientId)

		args := map[string]string{
			"${client}": id,
		}

		data, err := executeQuery(args, "client")

		if res, err := schemeConf.ConvertToClients(scheme, data); err == nil {
			if len(res) > 0 {
				order.Client = res[0]
			}
		} else {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		result = append(result, order)
	}
	return result, err
}

func setClientShipment(res []schemeConf.Shipment, executeQuery sqlI.Execute) ([]schemeConf.Shipment, error) {
	result := make([]schemeConf.Shipment, 0, len(res))

	rd, err := os.OpenFile("client.json", os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	scheme, err := schemeConf.CreateScheme(rd)
	if err != nil {
		return nil, err
	}

	for _, shipment := range res {

		id := fmt.Sprintf("0x%s", shipment.ClientId)

		args := map[string]string{
			"${client}": id,
		}

		data, err := executeQuery(args, "client")

		if res, err := schemeConf.ConvertToClients(scheme, data); err == nil {
			if len(res) > 0 {
				shipment.Client = res[0]
			}
		} else {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		result = append(result, shipment)
	}
	return result, err
}

func setOrderSegments(res []schemeConf.Order, executeQuery sqlI.Execute) ([]schemeConf.Order, error) {
	result := make([]schemeConf.Order, 0, len(res))

	for _, order := range res {

		id := fmt.Sprintf("0x%s", order.Ref)

		args := map[string]string{
			"${order}": id,
		}

		if data, err := executeQuery(args, "order_segments"); err == nil {
			order.LoadSegments(data)
		} else {
			return nil, err
		}

		result = append(result, order)

	}
	return result, nil
}

func setShipmentSegments(res []schemeConf.Shipment, executeQuery sqlI.Execute) ([]schemeConf.Shipment, error) {
	result := make([]schemeConf.Shipment, 0, len(res))

	for _, shipment := range res {

		id := fmt.Sprintf("0x%s", shipment.Ref)

		args := map[string]string{
			"${shipment}": id,
		}

		if data, err := executeQuery(args, "shipment_segments"); err == nil {
			shipment.LoadSegments(data)
		} else {
			return nil, err
		}

		result = append(result, shipment)

	}
	return result, nil
}

func setReceptionPropertyes(res []schemeConf.Reception, executeQuery sqlI.Execute) ([]schemeConf.Reception, error) {

	result := make([]schemeConf.Reception, 0, len(res))

	rd, err := os.OpenFile("reception_propertyes.json", os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	scheme, err := schemeConf.CreateScheme(rd)
	if err != nil {
		return nil, err
	}

	for _, reception := range res {

		args := map[string]string{
			"${reception}": fmt.Sprintf("0x%s", reception.Id),
		}
		data, err := executeQuery(args, "reception_propertyes")

		if res, err := schemeConf.ConvertToAdditionalFields(scheme, data); err == nil {
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
