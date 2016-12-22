package cmd

import (
	"fmt"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/secret"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize valec enviornment",
	Long: `Initialize valec environment

These resources will be created:
  - KMS key and alias
  - DynamoDB table`,
	RunE: doInit,
}

func doInit(cmd *cobra.Command, args []string) error {
	keyExists, err := aws.KMS.KeyExists(secret.DefaultKMSKey)
	if err != nil {
		return errors.Wrap(err, "Failed to check existence of key alias.")
	}

	if keyExists {
		fmt.Printf("Key %q alreadly exists.\n", secret.DefaultKMSKey)
	} else {
		keyID, err := aws.KMS.CreateKey()
		if err != nil {
			return errors.Wrap(err, "Failed to create new key.")
		}

		if err := aws.KMS.CreateKeyAlias(keyID, secret.DefaultKMSKey); err != nil {
			return errors.Wrap(err, "Failed to attach alias to key.")
		}

		fmt.Printf("Key %s successfully created!\n", secret.DefaultKMSKey)
	}

	tableExists, err := aws.DynamoDB.TableExists(rootOpts.tableName)
	if err != nil {
		return errors.Wrap(err, "Failed to check existence of DynamoDB table.")
	}

	if tableExists {
		fmt.Printf("DynamoDB table %q alreadly exists.\n", rootOpts.tableName)
	} else {
		if err := aws.DynamoDB.CreateTable(rootOpts.tableName); err != nil {
			return errors.Wrapf(err, "Failed to create DynamoDB table. table=%s", rootOpts.tableName)
		}

		fmt.Printf("DynamoDB table %s successfully created!\n", rootOpts.tableName)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(initCmd)
}
