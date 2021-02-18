// SeanZH shanghai 2019

package logh

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEngine(t *testing.T) {
	mux := Engine()
	ts := httptest.NewServer(mux)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/srv")
	if err != nil {
		t.Error("http get error", err)
	}
	defer resp.Body.Close()
	bt, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error("read all error:", err)
	}
	if string(bt) != "hello" {
		t.Error("we expect hello, we got", string(bt))
	}
}
