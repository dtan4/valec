package cmd

import (
	"fmt"

	"github.com/dtan4/valec/aws"
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
	keyExists, err := aws.KMS.KeyExists(keyAlias)
	if err != nil {
		return errors.Wrap(err, "Failed to check existence of key alias.")
	}

	if keyExists {
		fmt.Printf("Key %q alreadly exists.\n", keyAlias)
	} else {
		keyID, err := aws.KMS.CreateKey()
		if err != nil {
			return errors.Wrap(err, "Failed to create new key.")
		}

		if err := aws.KMS.CreateKeyAlias(keyID, keyAlias); err != nil {
			return errors.Wrap(err, "Failed to attach alias to key.")
		}

		fmt.Printf("Key %s successfully created!\n", keyAlias)
	}

	tableExists, err := aws.DynamoDB.TableExists(tableName)
	if err != nil {
		return errors.Wrap(err, "Failed to check existence of DynamoDB table.")
	}

	if tableExists {
		fmt.Printf("DynamoDB table %q alreadly exists.\n", tableName)
	} else {
		if err := aws.DynamoDB.CreateTable(tableName); err != nil {
			return errors.Wrapf(err, "Failed to create DynamoDB table. table=%s", tableName)
		}

		fmt.Printf("DynamoDB table %s successfully created!\n", tableName)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(initCmd)
}
