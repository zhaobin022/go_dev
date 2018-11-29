package main

import (
	"testing"
)

func TestAdd(t *testing.T) {
	r := add(5, 8)
	if r != 13 {
		t.Error("erro :", 6, r)
	}
	t.Log("test pass")
}
