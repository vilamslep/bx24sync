package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	bx24 "github.com/vi-la-muerto/bx24sync"
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
