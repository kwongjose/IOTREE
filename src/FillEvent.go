package main

import (
	"fmt"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		fmt.Printf("[%s] Message = %s \n", snsRecord.Timestamp, snsRecord.Message)
	}
}

func main() {
	lambda.Start(handler)
}
