package tools

import (
	"log"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
)


var CertPath string
var KeyPath string


type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("request header is %+v", req.Header)
	w.Write([]byte("good"))
}

func RunHTTP3() {
	var h Handler
	http.Handle("/", &h)
	http3.ListenAndServeQUIC(":6666", CertPath, KeyPath, nil)
}
