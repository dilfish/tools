package tools

import (
	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type Http struct {
	t *testing.T
}

func (h *Http) Hello(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
	assert.Equal(h.t, "hello", r.Body.String())
	assert.Equal(h.t, http.StatusOK, r.Code)
}

func (h *Http) Status(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
	assert.Equal(h.t, 404, r.Code)
}

func (h *Http) Header(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
	assert.Equal(h.t, "test-header", rq.Header.Get("X-Header-Test"))
}

func TestNewLogMux(t *testing.T) {
	lm := NewLogMux("testdata/http.log", "test_")
	if lm == nil {
		t.Error("lm is nil")
	}
	lm.GET("/abc", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	lm.GET("/header", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	r := gofight.New()
	var h Http
	h.t = t
	r.GET("/abc").SetDebug(true).Run(lm, h.Hello)
	r.GET("/status").SetDebug(true).Run(lm, h.Status)
	r.GET("/").SetDebug(true).SetHeader(gofight.H{"X-Header-Test": "test-header"}).Run(lm, h.Header)
}
