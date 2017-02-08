package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

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
	RunE: doList,
}

var listOpts = struct {
	secretFile string
	showValues bool
}{}

func doList(cmd *cobra.Command, args []string) error {
	var (
		secrets []*secret.Secret
		err     error
	)

	if listOpts.secretFile == "" {
		if len(args) != 1 {
			return errors.New("Please specify namespace or secret file (-f FILE).")
		}
		namespace := args[0]

		secrets, err = aws.DynamoDB.ListSecrets(rootOpts.tableName, namespace)
		if err != nil {
			return errors.Wrapf(err, "Failed to load secrets from DynamoDB. namespace=%s", namespace)
		}

		if len(secrets) == 0 {
			return errors.Errorf("Namespace %s does not exist.", namespace)
		}
	} else {
		_, secrets, err = secret.LoadFromYAML(listOpts.secretFile)
		if err != nil {
			return errors.Wrapf(err, "Failed to load secrets from file. filename=%s", listOpts.secretFile)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	for _, secret := range secrets {
		plainValue, err := aws.KMS.DecryptBase64(secret.Key, secret.Value)
		if err != nil {
			return errors.Wrapf(err, "Failed to decrypt value. key=%q, value=%q", secret.Key, secret.Value)
		}

		if listOpts.showValues {
			fmt.Fprintf(w, "%s\t%s\n", secret.Key+":", plainValue)
		} else {
			fmt.Fprintln(w, secret.Key)
		}
	}

	w.Flush()

	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&listOpts.secretFile, "file", "f", "", "Secret file")
	listCmd.Flags().BoolVar(&listOpts.showValues, "show-values", false, "Show values")
}
