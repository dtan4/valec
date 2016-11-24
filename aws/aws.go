package aws

import (
	"github.com/dtan4/valec/aws/dynamodb"
	"github.com/dtan4/valec/aws/kms"
)

var (
	dynamoDBClient *dynamodb.Client
	kmsClient      *kms.Client
)

// DynampDB returns DynamoDB API Client and create new one if it does not exit
func DynamoDB() *dynamodb.Client {
	if dynamoDBClient == nil {
		dynamoDBClient = dynamodb.NewClient()
	}

	return dynamoDBClient
}

// KMS returns KMS API Client and create new one if it does not exist
func KMS() *kms.Client {
	if kmsClient == nil {
		kmsClient = kms.NewClient()
	}

	return kmsClient
}
