package tools

import (
	"testing"
)

type Config struct {
	C int `json:"c"`
}

func TestReadConfig(t *testing.T) {
	var conf Config
	err := ReadConfig("testdata/test.conf", &conf)
	if err != nil {
		t.Error("read config err", err)
	}
	if conf.C != 28 {
		t.Error("conf.c is not 28", conf.C)
	}
}

func TestRandInt(t *testing.T) {
	max := 50
	a := RandInt(max)
	b := RandInt(max)
	if int(a) > max || int(b) > max {
		t.Error("a is bigger than max", a, b)
	}
}

func TestRandStr(t *testing.T) {
	str := RandStr(2)
	if len(str) != 2 {
		t.Error("len is not 2", str)
	}
}

func TestReadFile(t *testing.T) {
	bt, err := ReadFile("testdata/t.file")
	if err != nil {
		t.Error("err is", err)
	}
	if string(bt) != "abc\n" {
		t.Error("abc is not", string(bt))
		panic(string(bt))
	}
}

func TestFileMd5(t *testing.T) {
	n, md5, err := FileMd5("testdata/file.md5")
	if err != nil {
		t.Error("err is", err)
	}
	if n != 5 || md5 != "e7df7cd2ca07f4f1ab415d457a6e1c13" {
		t.Error("n is, md5", n, md5)
	}
}
