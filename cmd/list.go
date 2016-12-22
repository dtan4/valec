package cmd

import (
	"fmt"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/secret"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [NAMESPACE]",
	Short: "List stored secrets",
	Long: `List stored secrets

To list secrets stored in DynamoDB, specify namespace:
  $ valec list NAMESPACE
to list secrets stored in local file, specify file:
  $ valec list -f qa.yaml

Encrypted values are decrypted and printed as plain text.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			secrets []*secret.Secret
			err     error
		)

		if secretFile == "" {
			if len(args) != 1 {
				return errors.New("Please specify namespace or secret file (-f FILE).")
			}
			namespace := args[0]

			secrets, err = aws.DynamoDB.ListSecrets(tableName, namespace)
			if err != nil {
				return errors.Wrapf(err, "Failed to load secrets from DynamoDB. namespace=%s", namespace)
			}

			if len(secrets) == 0 {
				return errors.Errorf("Namespace %s does not exist.", namespace)
			}
		} else {
			_, secrets, err = secret.LoadFromYAML(secretFile)
			if err != nil {
				return errors.Wrapf(err, "Failed to load secrets from file. filename=%s", secretFile)
			}
		}

		longestLength := longestKeyLength(secrets)

		for _, secret := range secrets {
			plainValue, err := aws.KMS.DecryptBase64(secret.Key, secret.Value)
			if err != nil {
				return errors.Wrapf(err, "Failed to decrypt value. key=%q, value=%q", secret.Key, secret.Value)
			}

			padding := ""
			for i := 0; i < longestLength-len(secret.Key); i++ {
				padding += " "
			}

			fmt.Printf("%s:%s %s\n", secret.Key, padding, plainValue)
		}

		return nil
	},
}

func longestKeyLength(secrets []*secret.Secret) int {
	longest := 0

	for _, secret := range secrets {
		if longest < len(secret.Key) {
			longest = len(secret.Key)
		}
	}

	return longest
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&secretFile, "file", "f", "", "Secret file")
}
