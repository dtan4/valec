package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/dtan4/valec/lib"
	"github.com/pkg/errors"
)

const (
	includeKey = ".include"
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
func (c *Client) Delete(table, namespace string, configs []*lib.Config) error {
	if len(configs) == 0 {
		return nil
	}

	writeRequests := []*dynamodb.WriteRequest{}

	var writeRequest *dynamodb.WriteRequest

	for _, config := range configs {
		writeRequest = &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					"namespace": &dynamodb.AttributeValue{
						S: aws.String(namespace),
					},
					"key": &dynamodb.AttributeValue{
						S: aws.String(config.Key),
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

// Insert creates / updates records of configs in DynamoDB table
func (c *Client) Insert(table, namespace string, configs []*lib.Config) error {
	if len(configs) == 0 {
		return nil
	}

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

	_, err := c.api.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: requestItems,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to insert items.")
	}

	return nil
}

var included = map[string]bool{}

// ListConfigs returns all configs in the given table and namespace
func (c *Client) ListConfigs(table, namespace string) ([]*lib.Config, error) {
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
		return []*lib.Config{}, errors.Wrapf(err, "Failed to list up configs. namespace=%s", namespace)
	}

	configs := []*lib.Config{}
	var config *lib.Config

	for _, item := range resp.Items {
		if *item["key"].S == includeKey {
			includeNamespace := *item["value"].S

			if included[includeNamespace] {
				return []*lib.Config{}, errors.Errorf("Circular includes has detected. sourceNamespace=%s, includeNamespace=%s", namespace, includeNamespace)
			}
			included[includeNamespace] = true

			includeConfigs, err := c.ListConfigs(table, includeNamespace)
			if err != nil {
				return []*lib.Config{}, errors.Wrapf(err, "Failed to include another namespace. sourceNamespace=%s, includeNamespace=%s", namespace, includeNamespace)
			}

			for _, config := range includeConfigs {
				configs = append(configs, config)
			}
		} else {
			config = &lib.Config{
				Key:   *item["key"].S,
				Value: *item["value"].S,
			}

			configs = append(configs, config)
		}
	}

	return configs, nil
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
