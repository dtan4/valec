package cmd

import (
	"fmt"

	"github.com/dtan4/valec/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get NAMESPACE KEY",
	Short: "Get secret",
	RunE:  doGet,
}

func doGet(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Please specify both namespace and key")
	}
	namespace, key := args[0], args[1]

	secret, err := aws.DynamoDB.Get(rootOpts.tableName, namespace, key)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve secret")
	}

	plainValue, err := aws.KMS.DecryptBase64(secret.Key, secret.Value)
	if err != nil {
		return errors.Wrap(err, "Failed to decrypt secret")
	}

	fmt.Printf(plainValue)

	return nil
}

func init() {
	RootCmd.AddCommand(getCmd)
}
