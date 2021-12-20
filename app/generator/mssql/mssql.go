package mssql

import (
	"database/sql"
	"net/url"
	"time"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/vi-la-muerto/bx24sync"
)

func GetDatabaseConnection(config bx24sync.DataBaseConnection) (db *sql.DB, err error) {
	connector, err := mssql.NewConnector(makeConnURL(config).String())

	if err != nil {
		return db, err
	}

	db = sql.OpenDB(connector)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	return db, err
}

func makeConnURL(config bx24sync.DataBaseConnection) *url.URL {
	return &url.URL{
		Scheme: "sqlserver",
		Host:   config.Socket.String(),
		User:   url.UserPassword(config.User, config.Password),
	}
}
