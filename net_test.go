package tools

import (
    "testing"
)


func TestIP2Num(t *testing.T) {
    n := IP2Num("1.1.1.1")
    if n != 16843009 {
        t.Error("expect 16843009, got", n)
    }
}


func TestNum2IP(t *testing.T) {
    n := Num2IP(16843009)
    if n != "1.1.1.1" {
        t.Error("expect 1.1.1.1 got", n)
    }
}
