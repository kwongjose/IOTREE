package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	TableName = "FillEvent"
)

// Data Item stuct for FillEvent
type fillData struct {
	Timestamp string
	FlowTime  float32
	FlowRate  float32
	TotalFlow float32
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	fmt.Println("STARTING")
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		fmt.Printf("[%s] Message = %s \n", snsRecord.Timestamp, snsRecord.Message)

		fillItem := makeFillItem(snsRecord.Message)
		dynamoItem := makeDynamoInput(fillItem)
		addDynamoItem(dynamoItem)
	}

	// TEMP CODE
	var FakeData = `{"Timestamp":"2021-01-02T15:04:05-07","FlowTime":30,"FlowRate":0.5,"TotalFlow":0.5}`

	fillItem := makeFillItem(FakeData)
	dynamoItem := makeDynamoInput(fillItem)
	addDynamoItem(dynamoItem)

}

// makeFillItem makes a new FillItem
func makeFillItem(jsonData string) fillData {
	var data fillData
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println(err)
	}

	return data
}

// makeItemInput returns the actual PutItem to add to DynamoDB table
func makeDynamoInput(data fillData) *dynamodb.PutItemInput {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		fmt.Println("Got error marshalling new movie item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	dynamoItem := dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(TableName),
	}
	return &dynamoItem
}

// addItem adds a new fill event item to dynamodb
func addDynamoItem(dynamoItem *dynamodb.PutItemInput) {
	// Create DynamoDB client
	svc := createDBClient()

	_, err := svc.PutItem(dynamoItem)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Successfully added '")
}

// createDBClient creates a DynamoDB client session
func createDBClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	return dynamodb.New(sess)
}
