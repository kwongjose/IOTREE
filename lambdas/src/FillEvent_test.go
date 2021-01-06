package main

import (
	"testing"
)

func TestMakeFillItem(t *testing.T) {
	item := makeFillItem()
	if item.FlowTime != 30 {
		t.Fatalf("FUCKING WRONG")
	}
}
