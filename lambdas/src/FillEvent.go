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
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	tableName = "FillEvent"
)

// Data Item stuct for FillEvent
type fillData struct {
	Timestamp string
	// Time in seconds
	FlowTime  float32
	FlowRate  float32
	TotalFlow float32
}

type DynamoInterface struct {
	Dynamo dynamodbiface.DynamoDBAPI
}

func main() {

	lambda.Start(handler)
}

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	dynamoIter := DynamoInterface{
		Dynamo: createDBClient(),
	}

	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		fmt.Printf("[%s] Message = %s \n", snsRecord.Timestamp, snsRecord.Message)

		// runDB(snsRecord.Message, dynamoIter.Dynamo)
		dynamoIter.runDB(snsRecord.Message)
	}
}

// runDB uses dependency injection
func runDB(message string, client dynamodbiface.DynamoDBAPI) error {
	fillItem := makeFillItem(message)
	dynamoItem := makeDynamoInput(fillItem)
	return addDynamoItem(dynamoItem, client)
}

// runDB here uses the interface directly
func (dyn DynamoInterface) runDB(message string) error {
	fillItem := makeFillItem(message)
	dynamoItem := makeDynamoInput(fillItem)
	return addDynamoItem(dynamoItem, dyn.Dynamo)
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
		TableName: aws.String(tableName),
	}
	return &dynamoItem
}

// addItem adds a new fill event item to dynamodb
func addDynamoItem(dynamoItem *dynamodb.PutItemInput, client dynamodbiface.DynamoDBAPI) error {
	// Create DynamoDB client
	// svc := createDBClient()

	_, err := client.PutItem(dynamoItem)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Success")
	return nil
}

// createDBClient creates a DynamoDB client session
func createDBClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	return dynamodb.New(sess)
}
