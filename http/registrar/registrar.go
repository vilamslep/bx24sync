package registrar

//TODO Logs to file

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
	"github.com/vi-la-muerto/bx24-service/kafka/producer"
	"github.com/vi-la-muerto/bx24-service/scheme"
)

//client
//order
//reception
//shipment

type Registrar struct {
	producer.Producer
	*http.Server
	Config scheme.Registrar
	CheckInputEvent
	log.Logger
}

type CheckInputEvent func([]byte) (bool, error)

func (f CheckInputEvent) Execute(body []byte) (bool, error) {
	return f(body)
}

func NewRegistrar(config scheme.Registrar) Registrar {
	s := Registrar{Config: config}
	s.setSettings()
	s.OpenWriter()

	return s
}

func (s *Registrar) setSettings() {
	s.Server = &http.Server{
		Addr:    s.Config.Endpoint.String(),
		Handler: handlers.LoggingHandler(os.Stdout, http.DefaultServeMux),
	}

	s.Producer = producer.Producer{
		Addr:      s.Config.Broker.String(),
		Topic:     s.Config.Topic,
		Partition: s.Config.Partition,
	}
}

func (s *Registrar) Run() {

	http.HandleFunc(fmt.Sprintf("/%s", s.Config.Method), s.handlerMainMethod(s.Config.MessageKey))

	log.Info("Start service")

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Can't start to listener: %s\n", err)
	}

}

func (s *Registrar) CheckInput(body []byte) bool {
	checkResult := true

	if s.CheckInputEvent != nil {
		if res, err := s.CheckInputEvent(body); err == nil {
			return res
		}
	}

	return checkResult
}

func (s *Registrar) handlerMainMethod(key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Permission denied"))
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {

			log.Errorf("getting body error: %s ", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Don't manage to get body"))
			return
		}

		log.Info(string(body))

		if !s.CheckInput(body) {

			log.Error("body isn't correctly")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Body isn't correctly"))
			return
		}

		if ok := s.WriteMessage(body, key); ok != nil {
			log.Errorf("writing to broker error: %s ", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Don't manage to write message"))
			return
		}

		w.Write([]byte("Message is writed"))

	}
}

func (s *Registrar) Close() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		s.CloseWriter()
		cancel()

	}()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v\n", err)
	}

	log.Info("Server Exited Properly\n")
}
