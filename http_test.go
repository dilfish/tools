package tools

import (
	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type Case struct {
	t *testing.T
}

func (c *Case) Hello(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
	assert.Equal(c.t, "hello", r.Body.String())
	assert.Equal(c.t, http.StatusOK, r.Code)
}

func (c *Case) Status(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
	assert.Equal(c.t, 404, r.Code)
}

func (c *Case) Header(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
	assert.Equal(c.t, "test-header", rq.Header.Get("X-Header-Test"))
}

func TestNewLogMux(t *testing.T) {
	lm := NewLogMux("testdata/http.log", "test_")
	if lm == nil {
		t.Error("lm is nil")
	}
	lm.Handle("/abc", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	lm.Handle("/header", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	r := gofight.New()
	var c Case
	c.t = t
	r.GET("/abc").SetDebug(true).Run(lm, c.Hello)
	r.GET("/status").SetDebug(true).Run(lm, c.Status)
	r.GET("/").SetDebug(true).SetHeader(gofight.H{"X-Header-Test": "test-header"}).Run(lm, c.Header)
}
