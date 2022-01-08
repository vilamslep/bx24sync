package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	bx24 "github.com/vi-la-muerto/bx24sync"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("execution error %v", err)
	}
}

func run() (err error) {

	config := bx24.NewRegistrarConfigFromEnv()

	if config.Topic == "" {
		return fmt.Errorf("not defined kafka topic. ")
	}

	writer := bx24.NewKafkaWriter(config.Broker.String(), config.Topic)

	router := bx24.NewRouter(os.Stdout, os.Stderr, true)

	settingRouter(router, config.CheckInput, writer)

	server := &http.Server{
		Addr:    config.Http.String(),
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("can't start to listener: %s\n", err)
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
		return fmt.Errorf("server shutdown failed:%+v", err)
	}
	return err
}

func settingRouter(r bx24.Router, enableCheckInput bool, kw *bx24.KafkaWriter) {

	var checkInputFunc bx24.CheckInput = nil

	if enableCheckInput {
		checkInputFunc = bx24.DefaultCheckInput
	}

	allowsMethods := []string{"POST"}

	r.AddMethod(
		bx24.NewHttpMethod(
			"/client", bx24.DefaultHandler(kw, "client"),
			checkInputFunc, allowsMethods))

	r.AddMethod(
		bx24.NewHttpMethod(
			"/order", bx24.DefaultHandler(kw, "order"),
			checkInputFunc, allowsMethods))

	r.AddMethod(
		bx24.NewHttpMethod(
			"/shipment", bx24.DefaultHandler(kw, "shipment"),
			checkInputFunc, allowsMethods))

	r.AddMethod(
		bx24.NewHttpMethod(
			"/reception", bx24.DefaultHandler(kw, "reception"),
			checkInputFunc, allowsMethods))
}
