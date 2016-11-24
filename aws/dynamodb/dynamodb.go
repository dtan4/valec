package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/dtan4/valec/lib"
	"github.com/pkg/errors"
)

// Client represents the wrapper of DynamoDB API client
type Client struct {
	client *dynamodb.DynamoDB
}

// NewClient creates new Client object
func NewClient() *Client {
	return &Client{
		client: dynamodb.New(session.New(), &aws.Config{}),
	}
}

// CreateTable creates new table for Valec
func (c *Client) CreateTable(table string) error {
	_, err := c.client.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String("namespace"),
				AttributeType: aws.String("S"),
			},
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String("key"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			&dynamodb.KeySchemaElement{
				AttributeName: aws.String("namespace"),
				KeyType:       aws.String("HASH"),
			},
			&dynamodb.KeySchemaElement{
				AttributeName: aws.String("key"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(table),
	})
	if err != nil {
		return errors.Wrapf(err, "Failed to create DynamoDB table. table=%s", table)
	}

	return nil
}

// Insert creates / updates records of configs in DynamoDB table
func (c *Client) Insert(table, namespace string, configs []*lib.Config) error {
	writeRequests := []*dynamodb.WriteRequest{}

	var writeRequest *dynamodb.WriteRequest

	for _, config := range configs {
		writeRequest = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"namespace": &dynamodb.AttributeValue{
						S: aws.String(namespace),
					},
					"key": &dynamodb.AttributeValue{
						S: aws.String(config.Key),
					},
					"value": &dynamodb.AttributeValue{
						S: aws.String(config.Value),
					},
				},
			},
		}
		writeRequests = append(writeRequests, writeRequest)
	}

	requestItems := make(map[string][]*dynamodb.WriteRequest)
	requestItems[table] = writeRequests

	_, err := c.client.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: requestItems,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to insert items.")
	}

	return nil
}
