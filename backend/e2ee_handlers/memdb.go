package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// ENV VARIABLES:
// - AWS_REGION
// - DDB_TABLE_CONN
type MemDBManager struct {
	client *dynamodb.DynamoDB
}

func NewMemDBManager() *MemDBManager {
	// Initialize DynamoDB client
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
		},
	))
	return &MemDBManager{
		client: dynamodb.New(sess),
	}
}

// Manager funcs
// get connection by connectionID
func (m *MemDBManager) GetConnection(connectionID string) (*ConnectionEntry, error) {
	// Query client
	result, err := m.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("DDB_TABLE_CONN")),
		Key: map[string]*dynamodb.AttributeValue{
			"connectionId": {
				S: aws.String(connectionID),
			},
		},
	})
	// If error, return
	if err != nil {
		return nil, err
	}
	// Cast to ConnectionEntry
	entry := &ConnectionEntry{}
	err = dynamodbattribute.UnmarshalMap(result.Item, entry)
	return entry, err
}

// retrieve connections by username
func (m *MemDBManager) GetUserConnections(username string) ([]*ConnectionEntry, error) {
	// Query DynamoDB
	result, err := m.client.Query(&dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("DDB_TABLE_CONN")),
		KeyConditionExpression: aws.String("username = :username"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(username),
			},
		},
	})
	// If error, return
	if err != nil {
		return nil, err
	}
	// Cast to ConnectionEntries
	entries := []*ConnectionEntry{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &entries)
	return entries, err
}
