package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mock struct {
	dynamodbiface.DynamoDBAPI
	dynResponse dynamodb.PutItemOutput
}

// PutItem mocks the PutItem func
func (mq mock) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &mq.dynResponse, nil
}

func TestInterface(t *testing.T) {
	dbMock := mock{
		DynamoDBAPI: nil,
		// Here we specify the details of the response for putItem
		dynResponse: dynamodb.PutItemOutput{},
	}

	// Create a item matching the interface to call runDB
	dynMock := DynamoInterface{
		Dynamo: dbMock,
	}

	err := dynMock.runDB(`{"Timestamp":"2021-01-02T15:04:05-07","FlowTime":30,"FlowRate":0.5,"TotalFlow":0.5}`)

	if err == nil {
		t.Logf("Good")
	}
}

func TestInjection(t *testing.T) {
	dbMock := mock{
		DynamoDBAPI: nil,
		dynResponse: dynamodb.PutItemOutput{},
	}

	err := runDB(`{"Timestamp":"2021-01-02T15:04:05-07","FlowTime":30,"FlowRate":0.5,"TotalFlow":0.5}`, dbMock)
	if err == nil {
		t.Logf("Good")
	}
}

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
