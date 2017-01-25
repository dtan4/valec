package dynamodb

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/dtan4/valec/aws/mock"
	"github.com/dtan4/valec/secret"
	"github.com/golang/mock/gomock"
)

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockDynamoDBAPI(ctrl)
	client := NewClient(api)

	if client.api != api {
		t.Errorf("client.api does not match.")
	}
}

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

func TestDelete_30items(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writeRequests1 := []*dynamodb.WriteRequest{}

	for i := 0; i < 25; i++ {
		writeRequests1 = append(writeRequests1, &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					"namespace": &dynamodb.AttributeValue{
						S: aws.String("test"),
					},
					"key": &dynamodb.AttributeValue{
						S: aws.String("BAZ" + strconv.Itoa(i)),
					},
				},
			},
		})
	}

	writeRequests2 := []*dynamodb.WriteRequest{}

	for i := 25; i < 30; i++ {
		writeRequests2 = append(writeRequests2, &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					"namespace": &dynamodb.AttributeValue{
						S: aws.String("test"),
					},
					"key": &dynamodb.AttributeValue{
						S: aws.String("BAZ" + strconv.Itoa(i)),
					},
				},
			},
		})
	}

	api := mock.NewMockDynamoDBAPI(ctrl)
	api.EXPECT().BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"valec": writeRequests1,
		},
	}).Return(&dynamodb.BatchWriteItemOutput{}, nil)
	api.EXPECT().BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"valec": writeRequests2,
		},
	}).Return(&dynamodb.BatchWriteItemOutput{}, nil)
	client := &Client{
		api: api,
	}

	secrets := []*secret.Secret{}

	for i := 0; i < 30; i++ {
		secrets = append(secrets, &secret.Secret{
			Key:   "BAZ" + strconv.Itoa(i),
			Value: strconv.Itoa(i),
		})
	}

	table := "valec"
	namespace := "test"
	if err := client.Delete(table, namespace, secrets); err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}
}

func TestDelete_nosecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockDynamoDBAPI(ctrl)
	client := &Client{
		api: api,
	}

	table := "valec"
	namespace := "test"
	secrets := []*secret.Secret{}

	if err := client.Delete(table, namespace, secrets); err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}
}

func TestDeleteNamespace(t *testing.T) {
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
					S: aws.String("BAR"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("fuga"),
				},
			},
		},
	}, nil)
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
				&dynamodb.WriteRequest{
					DeleteRequest: &dynamodb.DeleteRequest{
						Key: map[string]*dynamodb.AttributeValue{
							"namespace": &dynamodb.AttributeValue{
								S: aws.String("test"),
							},
							"key": &dynamodb.AttributeValue{
								S: aws.String("BAR"),
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

	table := "valec"
	namespace := "test"
	if err := client.DeleteNamespace(table, namespace); err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}
}

func TestDeleteNamespace_30items(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	items := []map[string]*dynamodb.AttributeValue{}

	for i := 0; i < 30; i++ {
		items = append(items, map[string]*dynamodb.AttributeValue{
			"namespace": &dynamodb.AttributeValue{
				S: aws.String("test"),
			},
			"key": &dynamodb.AttributeValue{
				S: aws.String("BAZ" + strconv.Itoa(i)),
			},
			"value": &dynamodb.AttributeValue{
				S: aws.String(strconv.Itoa(i)),
			},
		})
	}

	writeRequests1 := []*dynamodb.WriteRequest{}

	for i := 0; i < 25; i++ {
		writeRequests1 = append(writeRequests1, &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					"namespace": &dynamodb.AttributeValue{
						S: aws.String("test"),
					},
					"key": &dynamodb.AttributeValue{
						S: aws.String("BAZ" + strconv.Itoa(i)),
					},
				},
			},
		})
	}

	writeRequests2 := []*dynamodb.WriteRequest{}

	for i := 25; i < 30; i++ {
		writeRequests2 = append(writeRequests2, &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					"namespace": &dynamodb.AttributeValue{
						S: aws.String("test"),
					},
					"key": &dynamodb.AttributeValue{
						S: aws.String("BAZ" + strconv.Itoa(i)),
					},
				},
			},
		})
	}

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
		Items: items,
	}, nil)
	api.EXPECT().BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"valec": writeRequests1,
		},
	}).Return(&dynamodb.BatchWriteItemOutput{}, nil)
	api.EXPECT().BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"valec": writeRequests2,
		},
	}).Return(&dynamodb.BatchWriteItemOutput{}, nil)
	client := &Client{
		api: api,
	}

	table := "valec"
	namespace := "test"
	if err := client.DeleteNamespace(table, namespace); err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}
}

func TestGet(t *testing.T) {
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
			"key": &dynamodb.Condition{
				ComparisonOperator: aws.String(dynamodb.ComparisonOperatorEq),
				AttributeValueList: []*dynamodb.AttributeValue{
					&dynamodb.AttributeValue{
						S: aws.String("FOO"),
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
					S: aws.String("FOO"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("bar"),
				},
			},
		},
	}, nil)
	client := &Client{
		api: api,
	}

	expected := &secret.Secret{
		Key:   "FOO",
		Value: "bar",
	}

	table := "valec"
	namespace := "test"
	key := "FOO"
	actual, err := client.Get(table, namespace, key)
	if err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Secret does not match. expected: %v, actual: %v", expected, actual)
	}
}

func TestGet_nosecret(t *testing.T) {
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
			"key": &dynamodb.Condition{
				ComparisonOperator: aws.String(dynamodb.ComparisonOperatorEq),
				AttributeValueList: []*dynamodb.AttributeValue{
					&dynamodb.AttributeValue{
						S: aws.String("FOOBAR"),
					},
				},
			},
		},
	}).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{},
	}, nil)
	client := &Client{
		api: api,
	}

	table := "valec"
	namespace := "test"
	key := "FOOBAR"
	_, err := client.Get(table, namespace, key)
	if err == nil {
		t.Errorf("Error should be raised.")
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
								S: aws.String("BAR"),
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
			Key:   "BAR",
			Value: "fuga",
		},
	}

	table := "valec"
	namespace := "test"
	if err := client.Insert(table, namespace, secrets); err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}
}

func TestInsert_30items(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writeRequests1 := []*dynamodb.WriteRequest{}

	for i := 0; i < 25; i++ {
		writeRequests1 = append(writeRequests1, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"namespace": &dynamodb.AttributeValue{
						S: aws.String("test"),
					},
					"key": &dynamodb.AttributeValue{
						S: aws.String("BAZ" + strconv.Itoa(i)),
					},
					"value": &dynamodb.AttributeValue{
						S: aws.String(strconv.Itoa(i)),
					},
				},
			},
		})
	}

	writeRequests2 := []*dynamodb.WriteRequest{}

	for i := 25; i < 30; i++ {
		writeRequests2 = append(writeRequests2, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"namespace": &dynamodb.AttributeValue{
						S: aws.String("test"),
					},
					"key": &dynamodb.AttributeValue{
						S: aws.String("BAZ" + strconv.Itoa(i)),
					},
					"value": &dynamodb.AttributeValue{
						S: aws.String(strconv.Itoa(i)),
					},
				},
			},
		})
	}

	api := mock.NewMockDynamoDBAPI(ctrl)
	api.EXPECT().BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"valec": writeRequests1,
		},
	}).Return(&dynamodb.BatchWriteItemOutput{}, nil)
	api.EXPECT().BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"valec": writeRequests2,
		},
	}).Return(&dynamodb.BatchWriteItemOutput{}, nil)
	client := &Client{
		api: api,
	}

	secrets := []*secret.Secret{}

	for i := 0; i < 30; i++ {
		secrets = append(secrets, &secret.Secret{
			Key:   "BAZ" + strconv.Itoa(i),
			Value: strconv.Itoa(i),
		})
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
					S: aws.String("BAR"),
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
			Key:   "BAR",
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
					S: aws.String("BAR"),
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

func TestNamespaceExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockDynamoDBAPI(ctrl)

	client := &Client{
		api: api,
	}

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
					S: aws.String("BAR"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("fuga"),
				},
			},
		},
	}, nil)
	api.EXPECT().Query(&dynamodb.QueryInput{
		TableName: aws.String("valec"),
		KeyConditions: map[string]*dynamodb.Condition{
			"namespace": &dynamodb.Condition{
				ComparisonOperator: aws.String(dynamodb.ComparisonOperatorEq),
				AttributeValueList: []*dynamodb.AttributeValue{
					&dynamodb.AttributeValue{
						S: aws.String("foobar"),
					},
				},
			},
		},
	}).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{},
	}, nil)

	table := "valec"
	testcases := []struct {
		namespace string
		expected  bool
	}{
		{
			namespace: "test",
			expected:  true,
		},
		{
			namespace: "foobar",
			expected:  false,
		},
	}

	for _, tc := range testcases {
		actual, err := client.NamespaceExists(table, tc.namespace)
		if err != nil {
			t.Errorf("Error should not be raised. error: %s", err)
		}

		if actual != tc.expected {
			t.Errorf("Result does not match. namespace: %s, expected: %t", tc.namespace, tc.expected)
		}
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
