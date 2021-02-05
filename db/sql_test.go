// Copyright 2018 Sean.ZH

package db

import (
	"testing"
)

func TestInitDB(t *testing.T) {
	var conf DBConfig
	err := ReadConfig("testdata/mysql.conf", &conf)
	if err != nil {
		t.Error("read config error:", err)
	}
	db, err := InitDB(&conf)
	if err != nil {
		t.Error("db error", err)
	}
	defer db.Close()
	err = ReadConfig("testdata/fake.mysql.conf", &conf)
	if err != nil {
		t.Error("read fake conf error:", err)
	}
	_, err = InitDB(&conf)
	if err == nil {
		t.Error("we could link to an fake mysql")
	}
}
