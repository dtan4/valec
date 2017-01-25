package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	dotenvName       = ".env"
	dotenvSampleName = ".env.sample"
)

// dotenvCmd represents the dotenv command
var dotenvCmd = &cobra.Command{
	Use:   "dotenv NAMESPACE",
	Short: "Generate .env using .env.sample",
	RunE:  doDotenv,
}

var dotenvOpts = struct {
	quote bool
}{}

func doDotenv(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Please specify namespace")
	}
	namespace := args[0]

	secrets, err := aws.DynamoDB.ListSecrets(rootOpts.tableName, namespace)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve secrets")
	}

	if len(secrets) == 0 {
		return errors.Errorf("Namespace %s does not exist", namespace)
	}

	var (
		dotenv []string
		err2   error
	)

	if _, err := os.Stat(dotenvSampleName); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s does not exist. Dumping all secrets...\n", dotenvSampleName)

			dotenv, err2 = dumpAll(secrets, dotenvOpts.quote)
			if err2 != nil {
				return errors.Wrap(err, "Failed to dump all secrets")
			}
		} else {
			return errors.Wrapf(err, "Failed to get stat of dotenv template %s", dotenvSampleName)
		}
	} else {
		dotenv, err2 = dumpWithTemplate(secrets, dotenvOpts.quote, dotenvSampleName, false)
		if err2 != nil {
			return errors.Wrap(err, "Failed to dump secrets with dotenv template")
		}
	}

	body := []byte(strings.Join(dotenv, "\n") + "\n")
	if err := util.WriteFileWithoutSection(dotenvName, body); err != nil {
		return errors.Wrapf(err, "Failed to write dotenv file %s", dotenvName)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(dotenvCmd)

	dotenvCmd.Flags().BoolVarP(&dotenvOpts.quote, "quote", "q", false, "Quote values")
}
