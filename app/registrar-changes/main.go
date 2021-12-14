package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/vi-la-muerto/bx24sync"
)

func main() {

	router := bx24sync.NewRouter()

	router.AddMethod(bx24sync.HttpMethod{
		Path:         "/",
		Handler:      index,
		AllowMethods: []string{"GET", "POST"},
	})

	router.AddMethod(bx24sync.HttpMethod{
		Path:       "/client",
		Handler:    getHandlerClient("My msg"),
		CheckInput: chechInput,
	})

	s := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Can't start to listener: %s\n", err)
		}
	}()

}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadGateway)
	w.Write([]byte("Hello world"))
}

func getHandlerClient(msg string) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if body, err := io.ReadAll(req.Body); err == nil {
			content := string(body)
			fmt.Println(content)
		}
		w.Write([]byte(msg))
	})
}

func chechInput(r io.Reader) (res bool) {
	return res
}


// import (
// 	"os"
// 	"os/signal"
// 	"regexp"
// 	"strconv"
// 	"strings"
// 	"syscall"

// 	"github.com/vi-la-muerto/bx24-service/http/registrar"
// 	"github.com/vi-la-muerto/bx24-service/scheme"
// )

// //env
// const (
// 	kHost   = "KAFKA_HOST"
// 	kPort   = "KAFKA_PORT"
// 	kTopic  = "KAFKA_TOPIC"
// 	kPart   = "KAFKA_TOPIC_PARTITION"
// 	kKey    = "KAFKA_TOPIC_MESSAGE_KEY"
// 	sPort   = "SERVICE_PORT"
// 	sMethod = "SERVICE_REST_METHOD"
// )

// var (
// 	brokerAddr string
// 	brokerPort int
// 	topic      string
// 	partition  int
// 	messageKey string
// 	srvPort    int
// 	restMethod string
// )

// func main() {

// 	setSettingsFromEnv()

// 	config := scheme.Registrar{
// 		ProducerConfig: scheme.ProducerConfig{
// 			Broker:     scheme.Socket{Host: brokerAddr, Port: brokerPort},
// 			Topic:      topic,
// 			Partition:  partition,
// 			MessageKey: messageKey,
// 		},
// 		Endpoint: scheme.Endpoint{
// 			Socket: scheme.Socket{Host: "", Port: srvPort},
// 			Method: restMethod,
// 		},
// 	}

// 	s := registrar.NewRegistrar(config)
// 	s.CheckInputEvent = handlerCheckInput

// 	//new chanel and exec subscription for handing
// 	done := make(chan os.Signal, 1)
// 	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

// 	go s.Run()

// 	//wait signal
// 	<-done

// 	s.Close()
// }

// func setSettingsFromEnv() {
// 	brokerAddr = getEnvWithFallback(kHost, "172.19.0.3")
// 	brokerPort = stringToInt(os.Getenv(kPort), 9092)
// 	topic = getEnvWithFallback(kTopic, "changes")
// 	partition = stringToInt(os.Getenv(kPart), 0)
// 	messageKey = getEnvWithFallback(kKey, "client")
// 	srvPort = stringToInt(os.Getenv(sPort), 35671)
// 	restMethod = getEnvWithFallback(sMethod, "client")
// }

// func getEnvWithFallback(env string, fallback string) string {
// 	val := os.Getenv(env)
// 	if len(val) == 0 {
// 		return fallback
// 	}

// 	return val
// }

// func stringToInt(val string, fallback int) int {
// 	if res, err := strconv.Atoi(val); err != nil {
// 		return fallback
// 	} else {
// 		return res
// 	}
// }

// func handlerCheckInput(body []byte) (bool, error) {

// 	content := strings.ReplaceAll(string(body), "\n", "")

// 	regStr := `^{"#",+[[:xdigit:]]{8}(-[[:xdigit:]]{4}){3}-[[:xdigit:]]{12},[\d]{1,6}:[[:xdigit:]]{32}}$`

// 	return regexp.MatchString(regStr, content)
// }
