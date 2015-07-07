package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	key := "XXX"
	svc := dynamodb.New(&aws.Config{Region: "us-west-1", LogLevel: 0})
	fmt.Printf("%+v\n", svc)
	result, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Tables:")
	for _, table := range result.TableNames {
		log.Println(*table)
	}

	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{ // Required
			"HashKeyElement": {
				S: aws.String(key),
			},
		},
		TableName: aws.String("user_profiles"), // Required
		AttributesToGet: []*string{
			aws.String("user_key"),     // Required
			aws.String("version"),      // Required
			aws.String("user_profile"), // Required
			// More values...
		},
	}
	fmt.Println(awsutil.StringValue(params))
	resp, err := svc.GetItem(params)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Generic AWS error with Code, Message, and original error (if any)
			fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				// A service error occurred
				fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
			}
		} else {
			// This case should never be hit, the SDK should always return an
			// error which satisfies the awserr.Error interface.
			fmt.Println(err.Error())
		}
	}
	fmt.Printf("%+v\n", resp)

}
