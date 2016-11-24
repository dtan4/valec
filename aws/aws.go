package aws

import (
	"github.com/dtan4/valec/aws/kms"
)

var (
	kmsClient *kms.Client
)

// KMS returns KMS API Client and create new one if it does not exist
func KMS() *kms.Client {
	if kmsClient == nil {
		kmsClient = kms.NewClient()
	}

	return kmsClient
}
