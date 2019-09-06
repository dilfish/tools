package tools

import (
	"testing"
	"net/http"
	"time"
	"net/url"
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
}


func TestOpenReqLogDB(t *testing.T) {
	var conf MgoConfig
	err := ReadConfig("testdata/mongo.conf", &conf)
	if err != nil {
		t.Error("no such mgo config")
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
