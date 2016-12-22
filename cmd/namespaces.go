package cmd

import (
	"fmt"

	"github.com/dtan4/valec/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// namespacesCmd represents the namespaces command
var namespacesCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "List all namespaces",
	RunE:  doNamespaces,
}

func doNamespaces(cmd *cobra.Command, args []string) error {
	namespaces, err := aws.DynamoDB.ListNamespaces(tableName)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve namespaces.")
	}

	for _, namespace := range namespaces {
		fmt.Println(namespace)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(namespacesCmd)
}
