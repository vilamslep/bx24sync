package generator

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	bx24 "github.com/vi-la-muerto/bx24sync"
	"github.com/vi-la-muerto/bx24sync/app"
	"github.com/vi-la-muerto/bx24sync/app/generator/mssql"
)

func Run() (err error) {

	config := bx24.NewGeneratorConfigFromEnv()

	router := bx24.NewRouter(os.Stdout, os.Stderr, true)

	db, err := mssql.GetDatabaseConnection(config.DB)

	if err != nil {
		return err
	}
	settingRouter(router, config.CheckInput, db, config.StorageQueryTxt)

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
		checkInputFunc = app.DefaultCheckInput
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

	textsQuery := map[string]string{
		"client":               string(txtCl),
		"reception":            string(txtRecept),
		"reception_propertyes": string(txtReceptProp),
	}

	r.AddMethod(
		bx24.HttpMethod{
			Path:         "/client",
			Handler:      HandlerClientWithDatabaseConnection(bx24.ExecuteQuery(db, textsQuery)),
			CheckInput:   checkInputFunc,
			AllowMethods: allowsMethods,
		})

	r.AddMethod(
		bx24.HttpMethod{
			Path:         "/reception",
			Handler:      HandlerReceptionWithDatabaseConnection(bx24.ExecuteQuery(db, textsQuery)),
			CheckInput:   checkInputFunc,
			AllowMethods: allowsMethods,
		})

	r.AddMethod(
		bx24.HttpMethod{
			Path:         "/order",
			Handler:      HandlerOrderWithDatabaseConnection(bx24.ExecuteQuery(db, textsQuery)),
			CheckInput:   checkInputFunc,
			AllowMethods: allowsMethods,
		})

	r.AddMethod(
		bx24.HttpMethod{
			Path:         "/shipment",
			Handler:      HandlerShipmentWithDatabaseConnection(bx24.ExecuteQuery(db, textsQuery)),
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
