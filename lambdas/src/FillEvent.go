package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Data Item stuct for FillEvent
type FillItem struct {
	Timestamp string
	FlowTime  int
	FlowRate  float64
	TotalFlow float64
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	fmt.Println("STARTING")
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		fmt.Printf("[%s] Message = %s \n", snsRecord.Timestamp, snsRecord.Message)
	}

	fillItem := makeFillItem()
	dynamoItem := makeItemInput(fillItem)
	addItem(dynamoItem)

}

// makeFillItem makes a new FillItem
func makeFillItem() FillItem {
	// Assume Rate is L per min
	return FillItem{
		Timestamp: "2021-01-02T15:04:05-0700",
		FlowTime:  30,
		FlowRate:  0.5,
		TotalFlow: (.5 / 60) * 30,
	}
}

// makeItemInput returns the actual PutItem to add to DynamoDB table
func makeItemInput(fillItem FillItem) *dynamodb.PutItemInput {
	av, err := dynamodbattribute.MarshalMap(fillItem)
	if err != nil {
		fmt.Println("Got error marshalling new movie item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	tableName := "FillEvent"
	dynamoItem := dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	return &dynamoItem
}

// addItem adds a new fill event item to dynamodb
func addItem(dynamoItem *dynamodb.PutItemInput) {
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
