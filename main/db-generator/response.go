package main

import(
	"net/http"
)

func writeServerError(w http.ResponseWriter, err error) error {
	writeResponse(w, http.StatusInternalServerError, []byte(err.Error()))
	return err
}

func writeBadRequests(w http.ResponseWriter, err error, msg []byte) error {
	writeResponse(w, http.StatusBadRequest, msg)
	return err
}

func writeResponse(w http.ResponseWriter, status int, msg []byte) {
	w.WriteHeader(status)
	w.Write(msg)
}
