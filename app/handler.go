package app

import (
	"context"
	"io"
	"net/http"

	"github.com/segmentio/kafka-go"
)

//handler
func DefaultHandler(writer *kafka.Writer, key string) func(http.ResponseWriter, *http.Request) error {
	return func(rw http.ResponseWriter, req *http.Request) error {
		
		body, err := io.ReadAll(req.Body)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Don't manage to get body"))
			return err
		}

		msg := kafka.Message{
			Key: []byte(key),
			Value: body,
		}

		if err := writer.WriteMessages(context.Background(), msg); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Don't manager to write message. Try it later"))
		
			return err
		} 

		rw.Write([]byte("Message is writed"))
		return nil
	}
}