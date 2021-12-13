package generator

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
	"github.com/vi-la-muerto/bx24-service/scheme"
	"github.com/vi-la-muerto/bx24-service/scheme/bitrix24"
)

type Service struct {
	*http.Server
	*sql.DB
	QueryText map[string]string
}

func NewServer(config scheme.GeneratorConfig) Service {

	connector, err := mssql.NewConnector(makeConnURL(config.DB).String())

	if err != nil {
		log.Fatalf("creating connector: %s\n", err.Error())
	}

	db := sql.OpenDB(connector)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	s := Service{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.Web.Port),
			Handler: handlers.LoggingHandler(os.Stdout, http.DefaultServeMux),
		},
		DB: db,
	}

	s.QueryText = make(map[string]string)

	fpath := fmt.Sprintf("%s/client_main.sql", config.StorageQueryTxt)

	bContent, err := ioutil.ReadFile(fpath)

	if err != nil {
		log.Fatalf("reading file with sql query: %s\n", err.Error())
	}

	s.QueryText["main_client"] = string(bContent)

	http.HandleFunc("/client", s.handlerMainMethod())

	return s
}

func (s *Service) handlerMainMethod() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Permission denied"))
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {

			log.Errorf("Don't manage to get body: %s\n", err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Don't manage to get body"))
			return
		}

		log.Infof("Body: %s\n", string(body))

		if !checkInput(body) {

			log.Error("Body isn't correctly")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Body isn't correctly"))
			return
		}

		pathWithId := strings.Replace(strings.Split(string(body), ",")[2], "}", "", 1)

		id := fmt.Sprintf("0x%s", strings.ToUpper(strings.Split(pathWithId, ":")[1]))

		queryText := strings.ReplaceAll(s.QueryText["main_client"], "${client}", id)

		data, err := s.executeAndReadQuery(queryText)

		if err != nil {

			log.Errorf("getting from db: %s", err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		contact, err := transformToSctruct(data)

		if err != nil {

			log.Errorf("transforming contact: %s", err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))

			return
		}

		if clientJson, err := json.Marshal(contact); err == nil {

			w.Write(clientJson)

		} else {

			log.Error(err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

	}
}

func checkInput(body []byte) bool {
	content := strings.ReplaceAll(string(body), "\n", "")

	regStr := `^{"#",+[[:xdigit:]]{8}(-[[:xdigit:]]{4}){3}-[[:xdigit:]]{12},[\d]{1,6}:[[:xdigit:]]{32}}$`

	matched, err := regexp.MatchString(regStr, content)

	if err != nil {
		log.Fatalf("checking input: %s\n ", err.Error())
		return false
	}

	return matched

}

func (s *Service) executeAndReadQuery(text string) ([]map[string]string, error) {
	rows, err := s.DB.Query(text)

	data := make([]map[string]string, 0)
	pretty := make(map[string]string)

	if err != nil {
		log.Errorf("executing query: %s", err.Error())
		return data, err
	}

	cols, err := rows.Columns()
	if err != nil {
		log.Errorf("not correctly response from db: %s", err.Error())
		return data, err
	}

	results := make([]interface{}, len(cols))
	for i := range results {
		results[i] = new(interface{})
	}

	for rows.Next() {
		if err := rows.Scan(results[:]...); err != nil {
			log.Errorf("scaning response from db: %s", err.Error())
			return data, err
		}
		for i := range results {
			val := *results[i].(*interface{})
			var str string

			if val == nil {
				str = "NULL"
			} else {
				switch v := val.(type) {
				case []byte:
					str = string(v)
				default:
					str = fmt.Sprintf("%v", v)
				}
			}

			pretty[cols[i]] = str
		}
		data = append(data, pretty)

	}

	return data, nil
}

func transformToSctruct(data []map[string]string) (bitrix24.Contact, error) {
	contact := bitrix24.Contact{}

	for _, value := range data {
		contact.TransoftFromMap(value)

		return contact, nil
	}

	return contact, errors.New("slice was empty")
}

func (s *Service) Run() {
	log.Info("Start service")

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("stating service: %s\n", err)
	}
}

func (s *Service) Close() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	//extra handing
	defer func() {
		s.DB.Close()
		cancel()
	}()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Info("Server Exited Properly")
}

func makeConnURL(config scheme.DataBaseAuth) *url.URL {

	return &url.URL{
		Scheme: "sqlserver",
		Host:   config.Socket.String(),
		User:   url.UserPassword(config.User, config.Password),
	}
}
