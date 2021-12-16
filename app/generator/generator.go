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

const (
	dbHostKey       = "DB_HOST"
	dbPortKey       = "DB_PORT"
	dbUserKey       = "DB_USER"
	dbPasswordKey   = "DB_PASSWORD"
	hsHostKey       = "HTTP_HOST"
	hsPortKey       = "HTTP_PORT"
	hsAddCheckInput = "HTTP_ADD_CHECK_INPUT"

	dbHostStd = "localhost"
	dbPortStd = 1433
	hsHostStd = "localhost"
	hsPortStd = 8095
)

func Run() (err error) {

	config := getConfigFromEnv()

	router := bx24.NewRouter(os.Stdout, os.Stderr, true)

	enCheckInput := app.StringToBool(os.Getenv(hsAddCheckInput), false)

	db, err := mssql.GetDatabaseConnection(config.DB)

	if err != nil {
		return err
	}
	settingRouter(router, enCheckInput, db, config.StorageQueryTxt)

	server := &http.Server{
		Addr:    config.Web.String(),
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Can't start to listener: %s\n", err)
	}

	return err
}

func getConfigFromEnv() bx24.GeneratorConfig {
	return bx24.GeneratorConfig{
		DB: bx24.DataBaseConnection{
			Socket: bx24.Socket{
				Host: app.GetEnvWithFallback(dbHostKey, hsHostStd),
				Port: app.StringToInt(os.Getenv(dbPortKey), dbPortStd),
			},
			BasicAuth: bx24.BasicAuth{
				User:     app.GetEnvWithFallback(dbUserKey, ""),
				Password: app.GetEnvWithFallback(dbPasswordKey, ""),
			},
		},
		Web: bx24.Socket{
			Host: app.GetEnvWithFallback(hsHostKey, hsHostStd),
			Port: app.StringToInt(os.Getenv(hsPortKey), hsPortStd),
		},
		StorageQueryTxt: "./sql",
	}
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
