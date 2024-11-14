package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type AWSAPIManager struct {
	client *apigatewaymanagementapi.ApiGatewayManagementApi
}

func init() {
	// Validate environment variables
	if os.Getenv("API_ENDPOINT") == "" {
		log.Fatal("API_ENDPOINT environment variable is required")
		panic("API_ENDPOINT environment variable is required")
	}
}

func NewAWSAPIManager() *AWSAPIManager {
	// Initialize AWS API Manager
	client := apigatewaymanagementapi.New(GetAWSSession(),
		&aws.Config{Endpoint: aws.String(os.Getenv("API_ENDPOINT"))})
	return &AWSAPIManager{client: client}
}

func (m *AWSAPIManager) PostToConnection(connectionID string, data []byte) error {
	// Post to connection
	_, err := m.client.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: &connectionID,
		Data:         data,
	})
	return err
}
