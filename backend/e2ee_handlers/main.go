package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

// Websocket - EKS
// Deployed behind AWS API Gateway with Websocket support
func main() {
	// Initialize API server
	server := NewServer()
	// Initialize API handler
	router := NewAPIHandler(server)
	// Routing
	lambda.Start(router.HandleRequests)
}
