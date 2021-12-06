package registrar

//TODO Logs to file

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/vi-la-muerto/bx24-service/kafka/producer"
)

//client
//order
//reception
//shipment

type Service struct {
	producer.Producer
	*http.Server
	//AccessblyAddrs []string
}

func NewServer(sPort int, bAddr string, bPort int, topic string, partition int) Service {

	s := Service{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", sPort),
			Handler: handlers.LoggingHandler(os.Stdout, http.DefaultServeMux),
		},
		Producer: producer.Producer{
			BrokerAddr: bAddr,
			Port:       bPort,
			Topic:      topic,
			Partition:  partition,
		},
	}

	s.Writer = producer.CreateWriter(s.Producer)

	http.HandleFunc("/client", s.handlerClient())

	return s
}

func (s *Service) Run() {
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	log.Println("Start service")

}

func (s *Service) Echo() string {
	return "Echo"
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

		if ok := s.WriteMessage(content, "client"); ok != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Don't manage to write message"))
			return
		}
		w.Write([]byte("Message is writed"))

	}
}

func (s *Service) Close() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	//extra handing
	defer func() {
		s.CloseWriter()
		cancel()

	}()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Print("Server Exited Properly")
}
