package tools

import (
	"testing"
)

func TestInitDB(t *testing.T) {
	var conf DBConfig
	err := ReadConfig("testdata/mysql.conf", &conf)
	db, err := InitDB(&conf)
	if err != nil {
		t.Error("db error", err)
	}
	defer db.Close()
}
