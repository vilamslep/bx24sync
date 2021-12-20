package app

import (
	"io"
	"net/http"
	bx24 "github.com/vi-la-muerto/bx24sync"
)

//handler
func DefaultHandler(writer *bx24.KafkaWriter, key string) func(http.ResponseWriter, *http.Request) error {
	return func(rw http.ResponseWriter, req *http.Request) error {

		body, err := io.ReadAll(req.Body)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Don't manage to get body"))
			return err
		}

		msg := bx24.Message{
			Key:   []byte(key),
			Value: body,
		}

		if err := writer.Write(msg); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Don't manager to write message. Try it later"))

			return err
		}

		rw.Write([]byte("Message is writed"))
		return nil
	}
}
