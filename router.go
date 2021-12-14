package bx24sync

import (
	"bytes"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type CheckInput func(r io.Reader) bool

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

		if !r.checkInputEvent(method, &logger, req) {
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

func (r *Router) checkInputEvent(method *HttpMethod, w http.ResponseWriter, req *http.Request) (ok bool) {

	ok = true
	if method.CheckInput != nil {
		if ok, req.Body = method.checkInput(req.Body); !ok {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Body isn't correctly"))
			return
		}
	}

	return ok
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

//read data and save thet into new reader for reader have data in method's handler
func (m *HttpMethod) checkInput(r io.Reader) (res bool, reader io.ReadCloser) {

	res = m.CheckInput == nil

	var buf bytes.Buffer

	tee := io.TeeReader(r, &buf)

	if m.CheckInput != nil {
		res = m.CheckInput(tee)
	} else {
		io.ReadAll(tee)
	}

	reader = io.NopCloser(bytes.NewReader(buf.Bytes()))

	return res, reader
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

