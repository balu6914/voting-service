package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
)

var db *dynamodb.DynamoDB

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Endpoint:    aws.String("http://dynamodb-local:8000"),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", ""),
	})
	if err != nil {
		log.Fatalf("Failed to connect to DynamoDB: %v", err)
	}

	db = dynamodb.New(sess)
}

func createTables() {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String("Users"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err := db.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	input = &dynamodb.CreateTableInput{
		TableName: aws.String("Votes"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("UserID"),
				KeyType:       aws.String("HASH"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("UserID"),
				AttributeType: aws.String("S"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err = db.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	var vote Vote
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Votes"),
		Item: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(vote.UserID),
			},
			"Choice": {
				S: aws.String(vote.Choice),
			},
		},
	}

	_, err = db.PutItem(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Vote recorded"))
}

func handleResults(w http.ResponseWriter, r *http.Request) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Votes"),
	}

	result, err := db.Scan(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	votes := []Vote{}
	for _, item := range result.Items {
		vote := Vote{
			UserID: *item["UserID"].S,
			Choice: *item["Choice"].S,
		}
		votes = append(votes, vote)
	}

	results := calculateResults(votes)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

func calculateResults(votes []Vote) map[string]int {
	results := make(map[string]int)
	for _, vote := range votes {
		results[vote.Choice]++
	}
	return results
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/vote", handleVote).Methods("POST")
	r.HandleFunc("/results", handleResults).Methods("GET")

	createTables()

	fmt.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
