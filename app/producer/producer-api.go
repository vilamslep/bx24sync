package producer

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
	"github.com/vi-la-muerto/bx24sync/app"
)

func Run() (err error) {

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
