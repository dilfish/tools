// Copyright 2018 Sean.ZH

package tools

import (
	"sort"
	"testing"
	"fmt"
)

func TestStateCountSlice(t *testing.T) {
	var sc []StateCount
	sc = append(sc, StateCount{"a", 1})
	sc = append(sc, StateCount{"b", 2})
	sort.Sort(StateCountSlice(sc))
	if sc[0].Count != 2 || sc[1].Count != 1 {
		t.Error("count error")
	}
}


func TestIPv6Counter(t *testing.T) {
	ipv6c := NewIPv6Counter(true)
	if ipv6c == nil {
		t.Error("bad args")
	}
	fmt.Println("start to renew, please wait")
	err := ipv6c.Renew()
	if err != nil {
		t.Error("renew error:", err)
	}
	str := ipv6c.String()
	if len(str) < 3 {
		t.Error("get bad string", str)
	}
	str = ipv6c.RealString()
	if len(str) < 3 {
		t.Error("get bad real string", str)
	}
	states := ipv6c.Struct()
	if len(states) < 1 {
		t.Error("we have bad struct", states)
	}
}
