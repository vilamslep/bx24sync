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

	content, err := os.ReadFile(fmt.Sprintf("%s/client.sql", queryTextsDir))
	if err != nil {
		return err
	}

	r.AddMethod(
		bx24.HttpMethod{
			Path:         "/client",
			Handler:      HandlerWithDatabaseConnection(ExecuteQuery(db, string(content))),
			CheckInput:   checkInputFunc,
			AllowMethods: allowsMethods,
		})

	return err
}
