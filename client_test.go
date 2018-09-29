package tools

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const userId = "666"

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/user/"+userId {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{"name":"test", "email":"a@example.com"}`))
	}
}

func TestClientApi(t *testing.T) {
	var h Handler
	mock := httptest.NewServer(&h)
	defer mock.Close()
	client := New(mock.URL)
	_, err := client.GetUser(userId)
	assert.Nil(t, err)
}
