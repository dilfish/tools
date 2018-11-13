// Copyright 2018 Sean.ZH

package tools

import (
    "testing"
    "sort"
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


func TestNewIPv6Counter(t *testing.T) {
    ipv6c := NewIPv6Counter()
    err := ipv6c.Renew()
    if err != nil {
        t.Error("renew", err)
    }
    str := ipv6c.String()
    if len(str) < 100 {
        t.Error("ipv6c string", str)
    }
}
