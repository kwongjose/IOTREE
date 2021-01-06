package main

import (
	"testing"
)

func TestMakeFillItem(t *testing.T) {
	var FakeData = `{"Timestamp":"2021-01-02T15:04:05-07","FlowTime":30,"FlowRate":0.5,"TotalFlow":0.5}`
	var testFail = false
	item := makeFillItem(FakeData)
	if item.FlowTime != 30 {
		t.Logf("FlowTime incorrect")
		testFail = true
	}
	if item.Timestamp != "2021-01-02T15:04:05-07" {
		t.Logf("Timestamp incorrect")
		testFail = true
	}
	if item.FlowRate != .5 {
		t.Logf("FlowRate incorrect")
		testFail = true
	}
	if item.TotalFlow != .5 {
		t.Logf("FlowTime incorrect")
		testFail = true
	}
	if testFail {
		t.Fail()
	}
}
