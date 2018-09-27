package tools

import (
	"testing"
)

func TestNewAppender(t *testing.T) {
	as, err := NewAppender("./test.log")
	if err != nil {
		t.Error("expect nil, got", err)
	}
	defer as.Close()
}

func TestInitLog(t *testing.T) {
	logger := InitLog("./test.log", "test_")
	if logger == nil {
		t.Errorf("logger is nil")
	}
}
