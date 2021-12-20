package changes

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	bx24 "github.com/vi-la-muerto/bx24sync"
	"github.com/vi-la-muerto/bx24sync/app"
)

const (
	kHostKey        = "KAFKA_HOST"
	kPortKey        = "KAFKA_PORT"
	kTopicKey       = "KAFKA_TOPIC"
	hsHostKey       = "HTTP_HOST"
	hsPortKey       = "HTTP_PORT"
	hsAddCheckInput = "HTTP_ADD_CHECK_INPUT"

	kHostStd  = "kafka"
	kPortStd  = 9092
	kTopicStd = "changes"
	hsHostStd = "localhost"
	hsPortStd = 8082
)

func Run() (err error) {

	config := getConfigFromEnv()

	writer := app.NewKafkaWriter(config.Broker.String(), config.Topic)

	router := bx24.NewRouter(os.Stdout, os.Stderr, true)

	enCheckInput := app.StringToBool(os.Getenv(hsAddCheckInput), false)

	settingRouter(router, enCheckInput, writer)

	server := &http.Server{
		Addr:    config.Http.String(),
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Start server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Can't start to listener: %s\n", err)
		}
	}()

	//wait signal
	<-done

	//free up resources
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		writer.Close()
		cancel()

	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed:%+v\n", err)
	}

	log.Println("server exited properly")

	return err
}

func getConfigFromEnv() bx24.RegistrarConfig {
	return bx24.RegistrarConfig{
		Http: bx24.Socket{
			Host: app.GetEnvWithFallback(hsHostKey, hsHostStd),
			Port: app.StringToInt(os.Getenv(hsPortKey), hsPortStd),
		},
		ProducerConfig: bx24.ProducerConfig{
			Broker: bx24.Socket{
				Host: app.GetEnvWithFallback(kHostKey, kHostStd),
				Port: app.StringToInt(os.Getenv(kPortKey), kPortStd),
			},
			Topic: app.GetEnvWithFallback(kTopicKey, kTopicStd),
		},
	}
}

func settingRouter(r bx24.Router, enableCheckInput bool, kw *app.KafkaWriter) {

	var checkInputFunc bx24.CheckInput = nil

	if enableCheckInput {
		checkInputFunc = app.DefaultCheckInput
	}

	allowsMethods := []string{"POST"}

	r.AddMethod(
		bx24.NewHttpMethod(
			"/client", app.DefaultHandler(kw, "client"),
			checkInputFunc, allowsMethods))

	r.AddMethod(
		bx24.NewHttpMethod(
			"/order", app.DefaultHandler(kw, "client"),
			checkInputFunc, allowsMethods))

	r.AddMethod(
		bx24.NewHttpMethod(
			"/shipment", app.DefaultHandler(kw, "client"),
			checkInputFunc, allowsMethods))

	r.AddMethod(
		bx24.NewHttpMethod(
			"/reseption", app.DefaultHandler(kw, "client"),
			checkInputFunc, allowsMethods))
}
