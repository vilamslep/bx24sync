package bx24sync

import (
	"bytes"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type CheckInput func(r io.Reader) (bool, error)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type Router struct {
	methods   map[string]*HttpMethod
	needLogBody bool
	accessLog *log.Logger
	errorLog  *log.Logger
}

func NewRouter(accessLog io.Writer, errorLog io.Writer, enableLoggingBody bool) (r Router) {
	r.methods = make(map[string]*HttpMethod)

	formatter := &log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}

	r.accessLog = log.New()
	r.accessLog.Out = accessLog
	r.accessLog.SetFormatter(formatter)
	
	r.errorLog = log.New()
	r.errorLog.Out = errorLog
	r.errorLog.SetFormatter(formatter)

	r.needLogBody = enableLoggingBody

	return r
}

func (r *Router) AddMethod(method HttpMethod) {
	r.methods[method.Path] = &method
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	url := *req.URL
	logger := writerLogger{ResponseWriter: w, status: 200}

	if method, ok := r.methods[url.Path]; ok && method.isAllow(req.Method) {

		if r.needLogBody {
			r.addLogBody(req)
		}

		if !r.checkInputEvent(method, &logger, req) {
			return
		}

		if len(method.AllowMethods) != 0  {
			if err := method.Handler(&logger, req); err != nil {
				r.addLogError(log.Fields{
					"method": method.Path,
					"point":  "ServeHTTP",
					"error":  err,
				}, "method handler fail")
			}
		} else {
			logger.WriteHeader(http.StatusMethodNotAllowed)
		}

	} else {
		http.NotFound(&logger, req)
	}

	r.addStateRequest(logger.Status(), url.Path, req.RemoteAddr)
}

func (r *Router) addLogBody(req *http.Request) {

	if r.needLogBody {
		var buf bytes.Buffer
		tee := io.TeeReader(req.Body, &buf)

		if body,err := io.ReadAll(tee); err == nil {
			r.addLogInfo(log.Fields{
				"method": req.URL.Path,
				"type" : req.Method,
				"content" : string(body),
			}, "Body")
		} else {
			r.addLogError(log.Fields{
				"method": req.URL.Path,
				"type": req.Method,
				"point": "addLogBody",
				"error": err,
			},"Adding requests body to log")
		}

		req.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))

	}
}

func (r *Router) checkInputEvent(method *HttpMethod, w http.ResponseWriter, req *http.Request) (bool) {
	if method.CheckInput != nil {
		var buf bytes.Buffer
		tee := io.TeeReader(req.Body, &buf)

		if ok, err := method.checkInput(tee); err != nil {

			r.addLogError(log.Fields{
				"method": method.Path,
				"point":  "checkInputEvent",
				"error":  err,
			}, "checking input fails")

			w.WriteHeader(http.StatusInternalServerError)

			return false
		} else if !ok {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Body isn't correctly"))

			return false
		} else {
			req.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
		}
	}
	return true
}

func (r *Router) addStateRequest(status int, urlPath string, addr string) {

	fields := log.Fields{"path": urlPath, "addr": addr}

	if status == 200 {
		r.addLogInfo(fields, status)
	} else if status > 400 && status < 500 {
		r.addLogWarn(fields, status)
	} else if status > 500 {
		r.addLogError(fields, status)
	}
}

func (r *Router) addLogInfo(fields log.Fields, msg interface{}) {
	r.accessLog.WithFields(fields).Info(msg)
}

func (r *Router) addLogError(fields log.Fields, msg interface{}) {
	r.errorLog.WithFields(fields).Error(msg)
}

func (r *Router) addLogWarn(fields log.Fields, msg interface{}) {
	r.errorLog.WithFields(fields).Warn(msg)
}

type HttpMethod struct {
	Path    string
	Handler HandlerFunc
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

func (m *HttpMethod) checkInput(r io.Reader) (bool, error) {
	return m.CheckInput(r)
}

func NewHttpMethod(path string, handler HandlerFunc, checkInput CheckInput, allowsMethods []string) HttpMethod {
	return HttpMethod{
		Path: path,
		Handler: handler,
		CheckInput: checkInput,
		AllowMethods: allowsMethods,
	}
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
