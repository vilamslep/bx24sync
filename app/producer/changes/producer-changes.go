package changes

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	bx24 "github.com/vi-la-muerto/bx24sync"
	"github.com/vi-la-muerto/bx24sync/app"
)

const (
	kHostKey  = "KAFKA_HOST"
	kPortKey  = "KAFKA_PORT"
	kTopicKey = "KAFKA_TOPIC"
	hsHostKey = "HTTP_HOST"
	hsPortKey = "HTTP_PORT"

	kHostStd  = "kafka"
	kPortStd  = 9092
	kTopicStd = "changes"
	hsHostStd = "localhost"
	hsPortStd = 8082
)

func Run() (err error) {
	
	config := getConfigFromEnv()
	
	writer := app.GetKafkaWriter(config.Broker.String(), config.Topic)

	router := getRouter(os.Stdin, os.Stderr, true, writer)

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

func getRouter(accessLog io.Writer,errorLog io.Writer, enableLogBody bool, kw *kafka.Writer ) (*bx24.Router){
	r := bx24.NewRouter(os.Stdout, os.Stderr, true)

	r.AddMethod(bx24.HttpMethod{
		Path:       "/client",
		Handler:    app.DefaultHandler(kw, "client"),
		CheckInput: app.HandlerCheckInput,
		AllowMethods: []string{"POST"},
	})

	r.AddMethod(bx24.HttpMethod{
		Path:       "/order",
		Handler:    app.DefaultHandler(kw, "order"),
		CheckInput: app.HandlerCheckInput,
		AllowMethods: []string{"POST"},
	})

	r.AddMethod(bx24.HttpMethod{
		Path:       "/shipment",
		Handler:    app.DefaultHandler(kw, "shipment"),
		CheckInput: app.HandlerCheckInput,
		AllowMethods: []string{"POST"},
	})

	r.AddMethod(bx24.HttpMethod{
		Path:       "/reception",
		Handler:    app.DefaultHandler(kw, "reception"),
		CheckInput: app.HandlerCheckInput,
		AllowMethods: []string{"POST"},
	})

	return &r
}



