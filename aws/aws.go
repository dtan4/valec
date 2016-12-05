package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	dynamodbapi "github.com/aws/aws-sdk-go/service/dynamodb"
	kmsapi "github.com/aws/aws-sdk-go/service/kms"
	"github.com/dtan4/valec/aws/dynamodb"
	"github.com/dtan4/valec/aws/kms"
	"github.com/pkg/errors"
)

var (
	// DynamoDB represents DynamoDB API client
	DynamoDB *dynamodb.Client
	// KMS represents KMS API client
	KMS *kms.Client
)

// Initialize initializes AWS API clients
func Initialize(region string) error {
	var (
		sess *session.Session
		err  error
	)

	if region == "" {
		sess, err = session.NewSession()
		if err != nil {
			return errors.Wrap(err, "Failed to create new AWS session.")
		}
	} else {
		sess, err = session.NewSession(&aws.Config{Region: aws.String(region)})
		if err != nil {
			return errors.Wrap(err, "Failed to create new AWS session.")
		}
	}

	DynamoDB = dynamodb.NewClient(dynamodbapi.New(sess))
	KMS = kms.NewClient(kmsapi.New(sess))

	return nil
}
