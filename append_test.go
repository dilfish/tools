// Copyright 2018 Sean.ZH

package tools

import (
	"os"
	"testing"
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
	_, err = NewAppender("testdata/a/appender.log")
	if err == nil {
		t.Error("expect err, got", err)
	}
	SetLog()
	Daemon()
}
