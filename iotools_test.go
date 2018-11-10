// Copyright 2018 Sean.ZH

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


func cbTestReadLineArr(line string) error {
    if line != "hello go" {
        return ErrBadFmt
    }
    return nil
}


func TestReadLineArr(t *testing.T) {
    err := ReadLineArr("testdata/linedata.txt", cbTestReadLineArr, 2)
    if err != nil {
        t.Error("read line error", err)
    }
}


func TestReadLine(t *testing.T) {
    err := ReadLine("testdata/linedata.txt", cbTestReadLineArr)
    if err != nil {
        t.Error("read line", err)
    }
}


func TestUnixToBJ(t *testing.T) {
    tm := UnixToBJ(1539599662)
    str := tm.Format("2006-01-02 15:04:05 -0700")
    if str != "2018-10-15 18:34:22 +0800" {
        t.Error("unix to bj", str)
    }
}


func TestUnixToPacific(t *testing.T) {
    tm := UnixToUSPacific(1539599662)
    str := tm.Format("2006-01-02 15:04:05 -0700")
    if str != "2018-10-15 03:34:22 -0700" {
        t.Error("unix to us pacific is", str)
    }
}


func TestUnixToUTC(t *testing.T) {
    tm := UnixToUTC(1539599662)
    str := tm.Format("2006-01-02 15:04:05 -0700")
    if str != "2018-10-15 10:34:22 +0000" {
        t.Error("unix to utc is", str)
    }
}
