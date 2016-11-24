package kms

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/pkg/errors"
)

// Client represents the wrapper of KMS API client
type Client struct {
	client *kms.KMS
}

// NewClient creates new KMSClient object
func NewClient() *Client {
	return &Client{
		client: kms.New(session.New(), &aws.Config{}),
	}
}

// Decrypt decrypts the given bawse64-encoded cipher text
func (c *Client) DecryptBase64(key, cipherText string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to decode as base64 string. text=%q", cipherText)
	}

	resp, err := c.client.Decrypt(&kms.DecryptInput{
		CiphertextBlob: decoded,
		EncryptionContext: map[string]*string{
			"key": aws.String(key),
		},
	})
	if err != nil {
		return "", errors.Wrapf(err, "Failed to decrypt the given cipherText. key=%s, value=%q", key, cipherText)
	}

	return string(resp.Plaintext), nil
}