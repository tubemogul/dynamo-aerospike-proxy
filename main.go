package main

import (
	"fmt"
	. "github.com/aerospike/aerospike-client-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/get/{user_key}", FetchFromAS)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func FetchFromAS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_key := vars["user_key"]
	client, err := NewClient(os.Args[1], 3000)
	key, err := NewKey("memory", "demoset", user_key)
	if err != nil {
		log.Fatal(err)
	}
	rec, err := client.Get(nil, key)
	if err != nil {
		log.Fatal(err)
	} else {
		if rec != nil {
			w.Header().Set("X-Cache", "HIT")
			fmt.Fprintln(w, "got:", rec)
		} else {
			w.Header().Set("X-Cache", "MISS")
			getFromDynamo(user_key)
		}
	}
}

func getFromDynamo(user_key string) {
	svc := dynamodb.New(&aws.Config{Region: "us-west-1", LogLevel: 0})
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{ // Required
			"user_key": {
				S: aws.String(user_key),
			},
		},
		TableName: aws.String("user_profiles"), // Required
		AttributesToGet: []*string{
			aws.String("user_key"),     // Required
			aws.String("version"),      // Required
			aws.String("user_profile"), // Required
		},
	}
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
	} else {
		if resp.Item["user_profile"] != nil {
			client, err := NewClient(os.Args[1], 3000)
			key, err := NewKey("memory", "demoset", user_key)
			if err != nil {
				log.Fatal(err)
			}
			bin1 := NewBin("bin1", resp.Item["user_profile"].GoString())
			client.PutBins(nil, key, bin1)
		}
	}
}
