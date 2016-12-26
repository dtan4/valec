package dynamodb

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/dtan4/valec/aws/mock"
	"github.com/dtan4/valec/secret"
	"github.com/golang/mock/gomock"
)

func TestCreateTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockDynamoDBAPI(ctrl)

	api.EXPECT().CreateTable(&dynamodb.CreateTableInput{
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
		TableName: aws.String("valec"),
	}).Return(&dynamodb.CreateTableOutput{}, nil)
	client := &Client{
		api: api,
	}

	table := "valec"
	if err := client.CreateTable(table); err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockDynamoDBAPI(ctrl)

	api.EXPECT().BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"valec": []*dynamodb.WriteRequest{
				&dynamodb.WriteRequest{
					DeleteRequest: &dynamodb.DeleteRequest{
						Key: map[string]*dynamodb.AttributeValue{
							"namespace": &dynamodb.AttributeValue{
								S: aws.String("test"),
							},
							"key": &dynamodb.AttributeValue{
								S: aws.String("BAZ"),
							},
						},
					},
				},
				&dynamodb.WriteRequest{
					DeleteRequest: &dynamodb.DeleteRequest{
						Key: map[string]*dynamodb.AttributeValue{
							"namespace": &dynamodb.AttributeValue{
								S: aws.String("test"),
							},
							"key": &dynamodb.AttributeValue{
								S: aws.String("FOO"),
							},
						},
					},
				},
			},
		},
	}).Return(&dynamodb.BatchWriteItemOutput{}, nil)
	client := &Client{
		api: api,
	}

	secrets := []*secret.Secret{
		&secret.Secret{
			Key:   "BAZ",
			Value: "1",
		},
		&secret.Secret{
			Key:   "FOO",
			Value: "bar",
		},
	}

	table := "valec"
	namespace := "test"
	if err := client.Delete(table, namespace, secrets); err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}
}

func TestInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockDynamoDBAPI(ctrl)

	api.EXPECT().BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"valec": []*dynamodb.WriteRequest{
				&dynamodb.WriteRequest{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]*dynamodb.AttributeValue{
							"namespace": &dynamodb.AttributeValue{
								S: aws.String("test"),
							},
							"key": &dynamodb.AttributeValue{
								S: aws.String("BAZ"),
							},
							"value": &dynamodb.AttributeValue{
								S: aws.String("1"),
							},
						},
					},
				},
				&dynamodb.WriteRequest{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]*dynamodb.AttributeValue{
							"namespace": &dynamodb.AttributeValue{
								S: aws.String("test"),
							},
							"key": &dynamodb.AttributeValue{
								S: aws.String("FOO"),
							},
							"value": &dynamodb.AttributeValue{
								S: aws.String("bar"),
							},
						},
					},
				},
				&dynamodb.WriteRequest{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]*dynamodb.AttributeValue{
							"namespace": &dynamodb.AttributeValue{
								S: aws.String("test"),
							},
							"key": &dynamodb.AttributeValue{
								S: aws.String("FOO"),
							},
							"value": &dynamodb.AttributeValue{
								S: aws.String("fuga"),
							},
						},
					},
				},
			},
		},
	}).Return(&dynamodb.BatchWriteItemOutput{}, nil)
	client := &Client{
		api: api,
	}

	secrets := []*secret.Secret{
		&secret.Secret{
			Key:   "BAZ",
			Value: "1",
		},
		&secret.Secret{
			Key:   "FOO",
			Value: "bar",
		},
		&secret.Secret{
			Key:   "FOO",
			Value: "fuga",
		},
	}

	table := "valec"
	namespace := "test"
	if err := client.Insert(table, namespace, secrets); err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}
}

func TestListSecrets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockDynamoDBAPI(ctrl)

	api.EXPECT().Query(&dynamodb.QueryInput{
		TableName: aws.String("valec"),
		KeyConditions: map[string]*dynamodb.Condition{
			"namespace": &dynamodb.Condition{
				ComparisonOperator: aws.String(dynamodb.ComparisonOperatorEq),
				AttributeValueList: []*dynamodb.AttributeValue{
					&dynamodb.AttributeValue{
						S: aws.String("test"),
					},
				},
			},
		},
	}).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("BAZ"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("1"),
				},
			},
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("FOO"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("bar"),
				},
			},
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("FOO"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("fuga"),
				},
			},
		},
	}, nil)
	client := &Client{
		api: api,
	}

	expected := []*secret.Secret{
		&secret.Secret{
			Key:   "BAZ",
			Value: "1",
		},
		&secret.Secret{
			Key:   "FOO",
			Value: "bar",
		},
		&secret.Secret{
			Key:   "FOO",
			Value: "fuga",
		},
	}

	table := "valec"
	namespace := "test"
	actual, err := client.ListSecrets(table, namespace)
	if err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Secrets does not match. expected: %v, actual: %v", expected, actual)
	}
}

func TestListNamespaces(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockDynamoDBAPI(ctrl)

	api.EXPECT().Scan(&dynamodb.ScanInput{
		TableName: aws.String("valec"),
	}).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("BAZ"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("1"),
				},
			},
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test2"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("FOO"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("bar"),
				},
			},
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("FOO"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("fuga"),
				},
			},
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test3"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("FOO"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("fuga"),
				},
			},
		},
	}, nil)
	client := &Client{
		api: api,
	}

	expected := []string{
		"test",
		"test2",
		"test3",
	}

	table := "valec"
	actual, err := client.ListNamespaces(table)
	if err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Namespaces does not match. expected: %q, actual: %q", expected, actual)
	}
}

func TestTableExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockDynamoDBAPI(ctrl)

	for i := 0; i < 2; i++ {
		api.EXPECT().ListTables(&dynamodb.ListTablesInput{}).Return(&dynamodb.ListTablesOutput{
			TableNames: []*string{
				aws.String("valec"),
			},
		}, nil)
	}
	client := &Client{
		api: api,
	}

	testcases := []struct {
		table    string
		expected bool
	}{
		{
			table:    "valec",
			expected: true,
		},
		{
			table:    "foobar",
			expected: false,
		},
	}

	for _, tc := range testcases {
		actual, err := client.TableExists(tc.table)
		if err != nil {
			t.Errorf("Error should not be raised. error: %s", err)
		}

		if actual != tc.expected {
			t.Errorf("Result does not match. table: %s, expected: %t", tc.table, tc.expected)
		}
	}
}
