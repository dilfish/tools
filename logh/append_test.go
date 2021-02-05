// Copyright 2018 Sean.ZH

package logh

import (
	"os"
	"testing"
)

func TestInitLog(t *testing.T) {
	t.Log("enter TestInitLog")
	defer t.Log("leave TestInitLog")
	logger := InitLog("testdata/initlog.log", "")
	if logger == nil {
		t.Error("expect logger, got nil")
	}
	logger = InitLog("", "")
	if logger != nil {
		t.Error("expect get nil, got", logger)
	}
}

func TestNewAppender(t *testing.T) {
	t.Log("enter TestNewAppender")
	defer t.Log("leave TestNewAppender")
	_, err := NewAppender("")
	if err == nil {
		t.Error("expect err got nil")
	}
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
	_, err = as.Write([]byte("end\n"))
	if err != nil {
		t.Error("wriet again err", err)
	}
	err = as.Restart()
	if err != nil {
		t.Error("restart error", err)
	}
	as.Close()
	for i := 0; i < 2; i++ {
		_, err = as.Write(nil)
		if err == nil {
			t.Error("expect not nil, got", err)
		}
	}
	SetLog()
	Daemon()
}
