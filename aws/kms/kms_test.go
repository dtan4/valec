package kms

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/dtan4/valec/aws/mock"
	"github.com/golang/mock/gomock"
)

func TestCreateKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockKMSAPI(ctrl)
	api.EXPECT().CreateKey(&kms.CreateKeyInput{
		Description: aws.String("Key for Valec"),
	}).Return(&kms.CreateKeyOutput{
		KeyMetadata: &kms.KeyMetadata{
			KeyId: aws.String("0123abcd-0123-9876-45ab-11aa22bb33cc"),
		},
	}, nil)
	client := &Client{
		api: api,
	}

	expected := "0123abcd-0123-9876-45ab-11aa22bb33cc"
	actual, err := client.CreateKey()
	if err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}

	if actual != expected {
		t.Errorf("KeyID does not match. expected: %q, actual: %q", expected, actual)
	}
}

func TestCreateKeyAlias(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockKMSAPI(ctrl)
	api.EXPECT().CreateAlias(&kms.CreateAliasInput{
		AliasName:   aws.String("alias/valec"),
		TargetKeyId: aws.String("0123abcd-0123-9876-45ab-11aa22bb33cc"),
	}).Return(&kms.CreateAliasOutput{}, nil)
	client := &Client{
		api: api,
	}

	keyID := "0123abcd-0123-9876-45ab-11aa22bb33cc"
	keyAlias := "valec"

	if err := client.CreateKeyAlias(keyID, keyAlias); err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}
}

func TestDecryptBase64(t *testing.T) {
	t.SkipNow()
}

func TestEncryptBase64(t *testing.T) {
	t.SkipNow()
}

func TestKeyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockKMSAPI(ctrl)
	for i := 0; i < 2; i++ {
		api.EXPECT().ListAliases(&kms.ListAliasesInput{}).Return(&kms.ListAliasesOutput{
			Aliases: []*kms.AliasListEntry{
				&kms.AliasListEntry{
					AliasName: aws.String("alias/valec"),
				},
				&kms.AliasListEntry{
					AliasName: aws.String("alias/valec-qa"),
				},
			},
		}, nil)
	}
	client := &Client{
		api: api,
	}

	testcases := []struct {
		keyAlias string
		expected bool
	}{
		{
			keyAlias: "valec",
			expected: true,
		},
		{
			keyAlias: "valec-prod",
			expected: false,
		},
	}

	for _, tc := range testcases {
		actual, err := client.KeyExists(tc.keyAlias)
		if err != nil {
			t.Errorf("Error should not be raised. error: %s", err)
		}

		if actual != tc.expected {
			t.Errorf("Result does not match. keyAlias: %s, expected: %t", tc.keyAlias, tc.expected)
		}
	}
}
