package cmd

import (
	"fmt"

	"github.com/dtan4/valec/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump NAMESPACE",
	Short: "Dump secrets in dotenv format",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Please specify namespace.")
		}
		namespace := args[0]

		configs, err := aws.DynamoDB().ListConfigs(tableName, namespace)
		if err != nil {
			return errors.Wrap(err, "Failed to retrieve configs.")
		}

		for _, config := range configs {
			plainValue, err := aws.KMS().DecryptBase64(config.Key, config.Value)
			if err != nil {
				return errors.Wrap(err, "Failed to decrypt value.")
			}

			fmt.Printf("%s=%s\n", config.Key, plainValue)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(dumpCmd)
}
