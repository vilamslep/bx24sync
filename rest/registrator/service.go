package registrator

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vi-la-muerto/bx24-service/kafka/producer"
)

//client
//order
//reception
//shipment

const (
// success          = 200
// badRequest       = 400
// notFound         = 404
// methodNotAllowed = 405
// serverError      = 500
)

type Service struct {
	Port int
	producer.Producer
	//AccessblyAddrs []string
}

func NewServer(port int) Service {

	s := Service{
		Port: port,
	}

	s.Producer = producer.Producer{
		BrokerAddr: "172.19.0.3",
		Port:       9092,
		Topic:      "clients.change",
		Partition:  0,
	}

	s.Writer = s.CreateWriter()

	http.HandleFunc("/client", handleClient)

	return s
}

func (s *Service) Start() {

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil))

}

func handleClient(w http.ResponseWriter, r *http.Request) {
	// if r.Method == "GET"{
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	w.Write([]byte("Method Not Allowed"))
	// 	return
	// }

	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func (s *Service) Close() {
	s.Producer.CloseWriter()
}
