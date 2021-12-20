package mssql

import (
	"database/sql"
	"time"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/vi-la-muerto/bx24sync"
)

func GetDatabaseConnection(config bx24sync.DataBaseConnection) (db *sql.DB, err error) {
	connector, err := mssql.NewConnector( config.MakeConnURL().String() )

	if err != nil {
		return db, err
	}

	db = sql.OpenDB(connector)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	return db, err
}


