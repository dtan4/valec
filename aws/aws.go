package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	dynamodbapi "github.com/aws/aws-sdk-go/service/dynamodb"
	kmsapi "github.com/aws/aws-sdk-go/service/kms"
	"github.com/dtan4/valec/aws/dynamodb"
	"github.com/dtan4/valec/aws/kms"
)

var (
	dynamoDBClient *dynamodb.Client
	kmsClient      *kms.Client
)

// DynamoDB returns DynamoDB API Client and create new one if it does not exit
func DynamoDB() *dynamodb.Client {
	api := dynamodbapi.New(session.New(), &aws.Config{})

	if dynamoDBClient == nil {
		dynamoDBClient = dynamodb.NewClient(api)
	}

	return dynamoDBClient
}

// KMS returns KMS API Client and create new one if it does not exist
func KMS() *kms.Client {
	api := kmsapi.New(session.New(), &aws.Config{})

	if kmsClient == nil {
		kmsClient = kms.NewClient(api)
	}

	return kmsClient
}
