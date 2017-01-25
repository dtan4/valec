package kms

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/pkg/errors"
)

// Client represents the wrapper of KMS API client
type Client struct {
	api kmsiface.KMSAPI
}

// NewClient creates new Client object
func NewClient(api kmsiface.KMSAPI) *Client {
	return &Client{
		api: api,
	}
}

// CreateKey creates new key and returns key ID
func (c *Client) CreateKey() (string, error) {
	resp, err := c.api.CreateKey(&kms.CreateKeyInput{
		Description: aws.String("Key for Valec"),
	})
	if err != nil {
		return "", errors.Wrap(err, "Failed to create key")
	}

	return *resp.KeyMetadata.KeyId, nil
}

// CreateKeyAlias attaches key alias to the given key
func (c *Client) CreateKeyAlias(keyID, keyAlias string) error {
	_, err := c.api.CreateAlias(&kms.CreateAliasInput{
		AliasName:   aws.String(keyAliasWithPrefix(keyAlias)),
		TargetKeyId: aws.String(keyID),
	})
	if err != nil {
		return errors.Wrapf(err, "Failed to create key alias %q for ID %q", keyAlias, keyID)
	}

	return nil
}

// DecryptBase64 decrypts the given base64-encoded cipher text
func (c *Client) DecryptBase64(key, cipherText string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to decode test %q as base64 string", cipherText)
	}

	resp, err := c.api.Decrypt(&kms.DecryptInput{
		CiphertextBlob: decoded,
		EncryptionContext: map[string]*string{
			"key": aws.String(key),
		},
	})
	if err != nil {
		return "", errors.Wrapf(err, "Failed to decrypt the given cipherText %q with key %q", cipherText, key)
	}

	return string(resp.Plaintext), nil
}

// EncryptBase64 encrypts the given text and return as base64-encoded cipher text
func (c *Client) EncryptBase64(keyAlias, key, text string) (string, error) {
	resp, err := c.api.Encrypt(&kms.EncryptInput{

		KeyId:     aws.String(keyAliasWithPrefix(keyAlias)),
		Plaintext: []byte(text),
		EncryptionContext: map[string]*string{
			"key": aws.String(key),
		},
	})
	if err != nil {
		return "", errors.Wrapf(err, "Failed to encrypt text with key %q", keyAlias)
	}

	return base64.StdEncoding.EncodeToString(resp.CiphertextBlob), nil
}

// KeyExists checks whether the given key exists or not
func (c *Client) KeyExists(keyAlias string) (bool, error) {
	resp, err := c.api.ListAliases(&kms.ListAliasesInput{})
	if err != nil {
		return false, errors.Wrap(err, "Failed to list key aliases")
	}

	aliasWithPrefix := keyAliasWithPrefix(keyAlias)

	for _, alias := range resp.Aliases {
		if *alias.AliasName == aliasWithPrefix {
			return true, nil
		}
	}

	return false, nil
}

func keyAliasWithPrefix(keyAlias string) string {
	// To use alias instead of KeyId, prefix 'alias/' is needed.
	// https://docs.aws.amazon.com/kms/latest/developerguide/programming-aliases.html
	return "alias/" + keyAlias
}
