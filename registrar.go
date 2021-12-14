package bx24sync

import (
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type CheckInput func(r io.Reader) bool

type HttpService struct {
	Server *http.Server
	Router
}

func NewHttpService() (hs HttpService) {

	hs.Router = NewRouter()

	hs.Server = &http.Server{
		Addr:    ":8082",
		Handler: hs.Router,
	}

	return hs
}

func (hs *HttpService) Run() {

	if err := hs.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Can't start to listener: %s\n", err)
	}
}

type HttpMethod struct {
	Path    string
	Handler http.HandlerFunc
	CheckInput
	AllowMethods []string
}

func (m *HttpMethod) isAllow(typeMethod string) bool {
	res := false
	for _, v := range m.AllowMethods {
		if v == typeMethod {
			res = true
			break
		}
	}
	return res
}
