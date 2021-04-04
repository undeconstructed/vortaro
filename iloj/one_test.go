package main

import "testing"

var CMPS = []struct {
	a string
	b string
	c bool
}{
	{"a", "aa", true},
	{"a", "b", true},
	{"a", "ĉ", true},
	{"ĉ", "a", false},
	{"ŝ", "z", true},
	{"b", "a", false},
	{"a'", "aa", true},
	{"-a", "a", true},
	{"-a", "aa", true},
}

func TestRootLess(t *testing.T) {
	for _, cmp := range CMPS {
		r := rootLess(cmp.a, cmp.b)
		if r != cmp.c {
			t.Errorf("%v", cmp)
		}
	}
}
