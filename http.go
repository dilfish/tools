package main

import (
	"io"
	"net/http"
)

type LogResponseWriter struct {
	ct     []byte
	status int
	w      http.ResponseWriter
}

type LogMux struct {
	mux    *http.ServeMux
	lw     *LogResponseWriter
	logger io.Writer
}

func NewLogMux(logger io.Writer) *LogMux {
	lm := &LogMux{}
	lm.mux = http.NewServeMux()
	lm.lw = &LogResponseWriter{}
	lm.logger = logger
	return lm
}

// implement http.ResponseWriter

func (lw *LogResponseWriter) Header() http.Header {
	return lw.w.Header()
}

func (lw *LogResponseWriter) Write(bt []byte) (int, error) {
	lw.ct = bt
	return lw.w.Write(bt)
}

func (lw *LogResponseWriter) WriteHeader(statusCode int) {
	lw.status = statusCode
	lw.w.WriteHeader(statusCode)
}

// implement http.ResponseWriter end

func (l *LogMux) Handle(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	l.mux.HandleFunc(pattern, handler)
}

func (l *LogMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, p := l.mux.Handler(r)
	l.lw.w = w
	h.ServeHTTP(l.lw, r)
	io.WriteString(l.logger, r.RequestURI+" "+p+" "+string(l.lw.ct))
}
