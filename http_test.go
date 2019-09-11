// Copyright 2018 Sean.ZH

package tools

import (
	"testing"
	"net/http"
	"encoding/json"
	"net/http/httptest"
)


func TestNewLogMux(t *testing.T) {
	lm := NewLogMux("", "")
	if lm != nil {
		t.Error("we expect nil, we get", lm)
	}
	lm = NewLogMux("testdata/mux.log", "test-")
	if lm == nil {
		t.Error("new log mux error")
	}
	lm.POST("/post", func (w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		w.WriteHeader(203)
		bt, _ := json.Marshal(h)
		w.Write(bt)
	})
	req := httptest.NewRequest("POST", "http://a.com/b", nil)
	w := httptest.NewRecorder()
	lm.ServeHTTP(w, req)
}
