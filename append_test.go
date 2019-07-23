// Copyright 2018 Sean.ZH

package tools

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestNewAppender(t *testing.T) {
	as, err := NewAppender("testdata/appender.log")
	if err != nil {
		t.Error("expect nil, got", err)
	}
	_, err = as.Write([]byte("start\n"))
	if err != nil {
		t.Error("expect nil, got", err)
	}
	err = os.Rename("testdata/appender.log", "testdata/appender.log.backup")
	if err != nil {
		t.Error("rename error", err)
	}
	err = syscall.Kill(as.pid, syscall.SIGUSR1)
	if err != nil {
		t.Error("send kill error", err)
	}
	time.Sleep(time.Second)
	_, err = as.Write([]byte("end\n"))
	if err != nil {
		t.Error("wriet again err", err)
	}
	as.Close()
	for i := 0; i < 2; i++ {
		_, err = as.Write(nil)
		if err == nil {
			t.Error("expect not nil, got", err)
		}
	}
	_, err = NewAppender("testdata/a/appender.log")
	if err == nil {
		t.Error("expect err, got", err)
	}
}

/*

func TestInitLog(t *testing.T) {
	logger := InitLog("testdata/log.log", "test_")
	if logger == nil {
		t.Error("logger is nil")
	}
	logger = InitLog("testdata/log.log", "")
	if logger == nil {
		t.Error("logger is nil 2")
	}
	logger = InitLog("testdata/a/b", "")
	if logger != nil {
		t.Error("logger is not nil 3")
	}
}

func TestDaemon(t *testing.T) {
	Daemon()
	if os.Stdin != nil {
		t.Error("stdin is not nil", os.Stdin)
	}
}
*/
