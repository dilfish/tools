// Copyright 2018 Sean.ZH

package tools

import (
	"sort"
	"testing"
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
