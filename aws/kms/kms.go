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

// NewClient creates new Client object
func NewClient() *Client {
	return &Client{
		client: kms.New(session.New(), &aws.Config{}),
	}
}

// DecryptBase64 decrypts the given base64-encoded cipher text
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
		return "", errors.Wrapf(err, "Failed to decrypt the given cipherText. key=%s, cipherText=%q", key, cipherText)
	}

	return string(resp.Plaintext), nil
}

// EncryptBase64 encrypts the given text and return as base64-encoded cipher text
func (c *Client) EncryptBase64(keyAlias, key, text string) (string, error) {
	resp, err := c.client.Encrypt(&kms.EncryptInput{
		// To use alias instead of KeyId, prefix 'alias/' is needed.
		// https://docs.aws.amazon.com/kms/latest/developerguide/programming-aliases.html
		KeyId:     aws.String("alias/" + keyAlias),
		Plaintext: []byte(text),
		EncryptionContext: map[string]*string{
			"key": aws.String(key),
		},
	})
	if err != nil {
		return "", errors.Wrapf(err, "Failed to encrypt text. keyAlias=%s, key=%q, text=%q", keyAlias, key, text)
	}

	return base64.StdEncoding.EncodeToString(resp.CiphertextBlob), nil
}
