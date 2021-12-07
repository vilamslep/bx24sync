package generator

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/vi-la-muerto/bx24-service/http/generator/scheme"
)

type Service struct {
	*http.Server
	*sql.DB
	QueryText map[string]string
}

func NewServer(config scheme.GeneratorConfig) Service {

	connector, err := mssql.NewConnector(makeConnURL(config.DB).String())

	if err != nil {
		log.Fatalln(err)
	}

	db := sql.OpenDB(connector)

	s := Service{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.Web.Port),
			Handler: handlers.LoggingHandler(os.Stdout, http.DefaultServeMux),
		},
		DB: db,
	}

	s.QueryText = make(map[string]string)

	fpath := fmt.Sprintf("%s/client_main.sql", config.QueryPath)

	bContent, err := ioutil.ReadFile(fpath)

	if err != nil {
		log.Fatalln(err)
	}

	s.QueryText["main_client"] = string(bContent)

	http.HandleFunc("/client", s.handlerClient())

	return s
}

func (s *Service) handlerClient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Permission denied"))
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Don't manage to get body"))
			return
		}

		content := strings.ReplaceAll(string(body), "\n", "")

		regStr := `^{"#",+[[:xdigit:]]{8}(-[[:xdigit:]]{4}){3}-[[:xdigit:]]{12},[\d]{1,6}:[[:xdigit:]]{32}}$`

		matched, err := regexp.MatchString(regStr, content)

		if err != nil {
			log.Fatal(err)
		}

		if !matched {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Body isn't correctly"))
			return
		}

		pathWithId := strings.Replace(strings.Split(content, ",")[2], "}", "", 1)

		id := fmt.Sprintf("0x%s", strings.ToUpper(strings.Split(pathWithId, ":")[1]))

		queryText := strings.ReplaceAll(s.QueryText["main_client"], "${client}", id)

		data := s.executeAndReadQuery(queryText)

		client, err := transformToSctruct(data)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))

			return
		}

		if clientJson, err := json.Marshal(client); err == nil {
			w.Write(clientJson)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

	}
}

func (s *Service) executeAndReadQuery(text string) []map[string]string {
	rows, err := s.DB.Query(text)

	if err != nil {
		log.Fatalln(err)
	}

	cols, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	data := make([]map[string]string, 0)
	pretty := make(map[string]string)

	results := make([]interface{}, len(cols))
	for i := range results {
		results[i] = new(interface{})
	}

	for rows.Next() {
		if err := rows.Scan(results[:]...); err != nil {
			panic(err)
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

	return data
}

func transformToSctruct(data []map[string]string) (scheme.Client, error) {
	client := scheme.Client{}

	for _, value := range data {
		client.TransoftFromMap(value)

		return client, nil
	}

	return client, errors.New("Slice was empty")
}

func (s *Service) Run() {
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	log.Println("Start service")
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

	log.Print("Server Exited Properly")
}

func makeConnURL(config scheme.DataBaseConfig) *url.URL {
	return &url.URL{
		Scheme: "sqlserver",
		Host:   config.Host + ":" + strconv.Itoa(config.Port),
		User:   url.UserPassword(config.User, config.Password),
	}
}
