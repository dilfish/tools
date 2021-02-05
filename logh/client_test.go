// Copyright 2018 Sean.ZH

package logh

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/abc" {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{"name":"test", "email":"a@example.com"}`))
	}
}

func TestClientApi(t *testing.T) {
	var h Handler
	mock := httptest.NewServer(&h)
	defer mock.Close()
	client := New(mock.URL, 1)
	_, err := client.Get("/abc")
	assert.Nil(t, err)
	client.SetBaseURL("a")
	a := client.GetBaseURL()
	if a != "a" {
		t.Error("get base url error", a)
	}
}
