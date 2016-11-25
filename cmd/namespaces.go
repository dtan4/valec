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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		namespaces, err := aws.DynamoDB().ListNamespaces(tableName)
		if err != nil {
			return errors.Wrap(err, "Failed to retrieve namespaces.")
		}

		for _, namespace := range namespaces {
			fmt.Println(namespace)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(namespacesCmd)
}
