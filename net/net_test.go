// Copyright 2018 Sean.ZH

package net

import (
	"testing"
)

func TestIP2Num(t *testing.T) {
	n := IP2Num("1.1.1.1")
	if n != 16843009 {
		t.Error("expect 16843009, got", n)
	}
	n = IP2Num("a")
	if n != 0 {
		t.Error("expect 0 but got", n)
	}
	n = IP2Num("1::1")
	if n != 0 {
		t.Error("expect 0 but got", n)
	}
}

func TestNum2IP(t *testing.T) {
	n := Num2IP(16843009)
	if n != "1.1.1.1" {
		t.Error("expect 1.1.1.1 got", n)
	}
}

func TestIPv62Num(t *testing.T) {
	inet, iint := IPv62Num("1::1")
	if inet != 281474976710656 || iint != 1 {
		t.Error("expect 2**48 and 1, got", inet, iint)
	}
	a, b := IPv62Num("1.1.1.1")
	if a != 0 || b != 0 {
		t.Error("we expect 0,0, but got", a, b)
	}
	a, b = IPv62Num("a")
	if a != 0 || b != 0 {
		t.Error("we expect 0,0, but got", a, b)
	}
}

func TestNum2IPv6(t *testing.T) {
	str := Num2IPv6(1, 1)
	if str != "0000:0000:0000:0001:0000:0000:0000:0001" {
		t.Error("expect 0000:0000:0000:0001:0000:0000:0000:0001, got", str)
	}
}

func TestDIG(t *testing.T) {
	_, err := DIG("baidu.com.", "114.114.114.114", "1.1.1.1")
	if err != nil {
		t.Error("expect nil, got", err)
	}
	_, err = DIG("baidu.com", "0", "1.1.1.1")
	if err == nil {
		t.Error("expect err but we got nil")
	}
}
