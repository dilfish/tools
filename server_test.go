// Copyright 2018 Sean.ZH

package tools

import (
	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type Server struct {
	t *testing.T
}

func (s *Server) Hello(r v2.HTTPResponse, rq v2.HTTPRequest) {
	assert.Equal(s.t, "hello", r.Body.String())
	assert.Equal(s.t, http.StatusOK, r.Code)
}

func TestHello(t *testing.T) {
	r := v2.New()
	e := Engine()
	var s Server
	s.t = t
	r.GET("/srv").SetDebug(true).Run(e, s.Hello)
}
