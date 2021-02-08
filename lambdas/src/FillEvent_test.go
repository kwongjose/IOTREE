package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mock struct {
	dynamodbiface.DynamoDBAPI
	dynResponse dynamodb.PutItemOutput
	insert      dynamodb.PutItemInput
}

// Store the value passed to PutItem
var insert dynamodb.PutItemInput

// PutItem mocks the PutItem func
func (mq mock) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {

	insert = *in
	return &mq.dynResponse, nil
}

func TestRunDB(t *testing.T) {
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

	if err != nil {
		t.Logf("FAILED %s: ", err)
		t.Fail()
	}

	if *insert.Item["FlowRate"].N != "0.5" {
		t.Logf("FlowRate incorrect")
		t.Fail()
	}

	if *insert.Item["TotalFlow"].N != "0.5" {
		t.Logf("TotalFlow incorrect")
		t.Fail()
	}
	if *insert.Item["FlowTime"].N != "30" {
		t.Logf("FlowTime incorrect")
		t.Fail()
	}
	if *insert.Item["Timestamp"].S != "2021-01-02T15:04:05-07" {
		t.Logf("Timestamp incorrect")
		t.Fail()
	}
	if *insert.TableName != "FillEvent" {
		t.Logf("TableName incorrect")
		t.Fail()
	}
}
