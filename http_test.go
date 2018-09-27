package tools

import (
	"testing"
)

func TestNewLogMux(t *testing.T) {
	lm := NewLogMux("./test.log", "test_")
	if lm == nil {
		t.Error("lm is nil")
	}
}
