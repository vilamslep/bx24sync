package registrator

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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
	producer.Producer
	*http.Server
	//AccessblyAddrs []string
}

func NewServer(port int) Service {

	s := Service{
		Server: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
		Producer: producer.Producer{
			BrokerAddr: "172.19.0.3",
			Port:       9092,
			Topic:      "clients.change",
			Partition:  0,
		},
	}

	s.Writer = producer.CreateWriter(s.Producer)

	http.HandleFunc("/client", handleClient)

	return s
}

func (s *Service) Run() {
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func handleClient(w http.ResponseWriter, r *http.Request) {
	// if r.Method == "GET"{
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	w.Write([]byte("Method Not Allowed"))
	// 	return
	// }


		
	// fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
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
