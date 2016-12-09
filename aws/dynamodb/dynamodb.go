package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/dtan4/valec/secret"
	"github.com/pkg/errors"
)

// Client represents the wrapper of DynamoDB API client
type Client struct {
	api dynamodbiface.DynamoDBAPI
}

// NewClient creates new Client object
func NewClient(api dynamodbiface.DynamoDBAPI) *Client {
	return &Client{
		api: api,
	}
}

// CreateTable creates new table for Valec
func (c *Client) CreateTable(table string) error {
	_, err := c.api.CreateTable(&dynamodb.CreateTableInput{
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

// Delete deletes records from DynamoDB table
func (c *Client) Delete(table, namespace string, secrets []*secret.Secret) error {
	if len(secrets) == 0 {
		return nil
	}

	writeRequests := []*dynamodb.WriteRequest{}

	var writeRequest *dynamodb.WriteRequest

	for _, secret := range secrets {
		writeRequest = &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					"namespace": &dynamodb.AttributeValue{
						S: aws.String(namespace),
					},
					"key": &dynamodb.AttributeValue{
						S: aws.String(secret.Key),
					},
				},
			},
		}
		writeRequests = append(writeRequests, writeRequest)
	}

	requestItems := make(map[string][]*dynamodb.WriteRequest)
	requestItems[table] = writeRequests

	_, err := c.api.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: requestItems,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to delete items.")
	}

	return nil
}

// Insert creates / updates records of secrets in DynamoDB table
func (c *Client) Insert(table, namespace string, secrets []*secret.Secret) error {
	if len(secrets) == 0 {
		return nil
	}

	writeRequests := []*dynamodb.WriteRequest{}

	var writeRequest *dynamodb.WriteRequest

	for _, secret := range secrets {
		writeRequest = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"namespace": &dynamodb.AttributeValue{
						S: aws.String(namespace),
					},
					"key": &dynamodb.AttributeValue{
						S: aws.String(secret.Key),
					},
					"value": &dynamodb.AttributeValue{
						S: aws.String(secret.Value),
					},
				},
			},
		}
		writeRequests = append(writeRequests, writeRequest)
	}

	requestItems := make(map[string][]*dynamodb.WriteRequest)
	requestItems[table] = writeRequests

	_, err := c.api.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: requestItems,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to insert items.")
	}

	return nil
}

// ListSecrets returns all secrets in the given table and namespace
func (c *Client) ListSecrets(table, namespace string) ([]*secret.Secret, error) {
	keyConditions := map[string]*dynamodb.Condition{
		"namespace": &dynamodb.Condition{
			ComparisonOperator: aws.String(dynamodb.ComparisonOperatorEq),
			AttributeValueList: []*dynamodb.AttributeValue{
				&dynamodb.AttributeValue{
					S: aws.String(namespace),
				},
			},
		},
	}
	params := &dynamodb.QueryInput{
		TableName:     aws.String(table),
		KeyConditions: keyConditions,
	}

	resp, err := c.api.Query(params)
	if err != nil {
		return []*secret.Secret{}, errors.Wrapf(err, "Failed to list up secrets. namespace=%s", namespace)
	}

	secrets := []*secret.Secret{}

	for _, item := range resp.Items {
		secret := &secret.Secret{
			Key:   *item["key"].S,
			Value: *item["value"].S,
		}

		secrets = append(secrets, secret)
	}

	return secrets, nil
}

// ListNamespaces returns all namespaces
func (c *Client) ListNamespaces(table string) ([]string, error) {
	resp, err := c.api.Scan(&dynamodb.ScanInput{
		TableName: aws.String(table),
	})
	if err != nil {
		return []string{}, errors.Wrapf(err, "Failed to retrieve items from DynamoDB table. table=%s", table)
	}

	nsmap := map[string]bool{}

	for _, item := range resp.Items {
		nsmap[*item["namespace"].S] = true
	}

	namespaces := []string{}

	for k := range nsmap {
		namespaces = append(namespaces, k)
	}

	return namespaces, nil
}

// TableExists check whether the given table exists or not
func (c *Client) TableExists(table string) (bool, error) {
	resp, err := c.api.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return false, errors.Wrap(err, "Failed to retrieve DynamoDB tables.")
	}

	for _, tableName := range resp.TableNames {
		if *tableName == table {
			return true, nil
		}
	}

	return false, nil
}
