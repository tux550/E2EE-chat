package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var myAWSSession *session.Session

func init() {
	// Validate environment variables
	if os.Getenv("AWS_REGION") == "" {
		log.Fatal("AWS_REGION environment variable is not set")
		panic("AWS_REGION environment variable is not set")
	}
	// Initialize AWS session
	myAWSSession = session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
		},
	))
}

func GetAWSSession() *session.Session {
	return myAWSSession
}
