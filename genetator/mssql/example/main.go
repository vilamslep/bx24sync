package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"strings"

	mssql "github.com/denisenkom/go-mssqldb"
)

const (
	server   = "proxy3"
	port     = 1433
	database = "dev"
	user     = "sa"
	password = "ssd445SQL"
	query    = "SELECT _IDRRef , _EnumOrder FROM dev.dbo._Enum563"
)

func main() {

	fContent, err := ioutil.ReadFile("clientbase.txt")

	if err != nil {
		log.Panic(err)
	}

	sContent := string(fContent)

	pathWithId := strings.Replace(strings.Split(sContent, ",")[2], "}", "", 1)

	id := fmt.Sprintf("0x%s", strings.ToUpper(strings.Split(pathWithId, ":")[1]))

	// log.Println(id)

	connString := makeConnURL().String()

	connector, err := mssql.NewConnector(connString)
	if err != nil {
		log.Println(err)
		return
	}

	query := fmt.Sprintf(`SELECT TOP 10 
	Client._IDRRef as Ref,
	Client._Description as Name,
	Client._Fld3795 as IsClient,
	Client._Marked as Deleted
	FROM dev.dbo._Reference151 Client WITH(NOLOCK)
	WHERE 
		Client._Marked = 0x00 
		AND Client._Fld3795 = 0x01 
		AND Client._IDRRef = %s
	`, id)

	db := sql.OpenDB(connector)
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var ref, name, isClient, deleted string

		rows.Scan(&ref, &name, &isClient, &deleted)

		fmt.Println(ref, name, isClient, deleted)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

}

func makeConnURL() *url.URL {
	return &url.URL{
		Scheme: "sqlserver",
		Host:   server + ":" + strconv.Itoa(port),
		User:   url.UserPassword(user, password),
	}
}
