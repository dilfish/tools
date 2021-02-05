package logh

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/dilfish/tools/io"
)

func TestRequestToInfo(t *testing.T) {
	var req http.Request
	tx := time.Now()
	req.Method = "POST"
	req.URL = &url.URL{}
	req.URL.Path = "/test"
	req.RemoteAddr = "1.1.1.1"
	ri := RequestToInfo(&req, tx)
	if ri.Time != tx {
		t.Error("requestinfo.t error", ri.Time, t)
	}
	if ri.ClientIP != "1.1.1.1" {
		t.Error("bad clientip", ri.ClientIP, "1.1.1.1")
	}
	if ri.Path != req.URL.Path {
		t.Error("bad path", ri.Path, req.URL.Path)
	}
	if ri.Method != req.Method {
		t.Error("bad method", ri.Method, req.Method)
	}
	req.RemoteAddr = "1.1.1.1:2222"
	ri = RequestToInfo(&req, tx)
	if ri.ClientIP != "1.1.1.1" {
		t.Error("bad ip:port", ri.ClientIP, "1.1.1.1")
	}
}

func TestNewRequestLogger(t *testing.T) {
	get := "/get"
	post := "/post"
	rl := NewRequestLogger(post, get)
	if rl.PostUrl != post || rl.GetUrl != get {
		t.Error("bad get/post", rl, get, post)
	}
	var ei ErrInfo
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bt, err := json.Marshal(ei)
		if err != nil {
			t.Error("marshal error:", err)
		}
		w.Write(bt)
	}))
	defer ts.Close()
	rl.PostUrl = ts.URL + "/post"
	req, err := http.NewRequest("POST", ts.URL+"/post", nil)
	if err != nil {
		t.Error("new request error:", err)
	}
	err = rl.PostOne(req)
	if err != nil {
		t.Error("post one error:", err)
	}
	rl.GetUrl = ts.URL + "/get"
	_, err = rl.GetStat(time.Now().Add(-time.Second), time.Now().Add(time.Second))
	if err != nil {
		t.Error("get stat", err)
	}
	rl.PostUrl = "/post"
	rl.GetUrl = "/get"
	err = rl.PostOne(req)
	if err == nil {
		t.Error("bad post one", err)
	}
	_, err = rl.GetStat(time.Now().Add(-time.Second), time.Now().Add(time.Second))
	if err == nil {
		t.Error("bad get one", err)
	}
}

func TestOpenReqLogDB(t *testing.T) {
	var conf MgoConfig
	err := io.ReadConfig("testdata/mongo.conf", &conf)
	if err != nil {
		t.Error("no such mgo config", err)
	}
	db := OpenReqLogDB(conf)
	if db == nil {
		t.Error("open mongo db error")
	}
	db.Close()
	conf.Username = "root"
	conf.Password = "ititititititiitiititititii"
	db = OpenReqLogDB(conf)
	if db != nil {
		t.Error("fake db open good:", db)
	}
}

func TestRequestLogger(t *testing.T) {
	var conf MgoConfig
	err := io.ReadConfig("testdata/mongo.conf", &conf)
	if err != nil {
		t.Error("read config error", err)
	}
	s := NewServeRequestLogger(conf)
	if s == nil {
		t.Error("new serve request error", err)
	}
	now := time.Now()
	var r RequestInfo
	r.Name = "test"
	r.Method = "POST"
	r.Path = "/do"
	r.ClientIP = "1.1.1.1"
	r.Time = now
	err = s.OneRequest(&r)
	if err != nil {
		t.Error("one request error", err)
	}
	_, err = s.GetStat(now.Add(-time.Second), now.Add(time.Second))
	if err != nil {
		t.Error("get state error:", err)
	}
}
