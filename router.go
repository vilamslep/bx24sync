package bx24sync

import (
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Router struct {
	methods map[string]*HttpMethod
}

func NewRouter() (r Router) {
	r.methods = make(map[string]*HttpMethod)

	return r
}

func (r *Router) AddMethod(method HttpMethod) {
	r.methods[method.Path] = &method
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	url := *req.URL

	logger := writerLogger{ResponseWriter: w, status: 200}

	if method, ok := r.methods[url.Path]; ok {

		ok := r.checkInput(*method, req.Body)
		if !ok {
			logger.WriteHeader(http.StatusBadGateway)
			logger.Write([]byte("Body isn't correctly"))

			return
		}

		if len(method.AllowMethods) == 0 || method.isAllow(req.Method) {
			method.Handler(&logger, req)
		} else {
			logger.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {
		http.NotFound(&logger, req)
	}

	log.WithFields(log.Fields{
		"Status":     logger.Status(),
		"Path":       url.Path,
		"RemoteAddr": req.RemoteAddr,
	}).Info("Access")
}

func (s *Router) checkInput(m HttpMethod, r io.Reader) (res bool) {
	
	res = m.CheckInput == nil
	
	if m.CheckInput != nil {
		res = m.CheckInput(r)
	}

	return res
}



type writerLogger struct {
	http.ResponseWriter
	status int
}

func (l *writerLogger) WriteHeader(code int) {
	l.status = code
	l.ResponseWriter.WriteHeader(code)
}

func (l *writerLogger) Status() int {
	return l.status
}
