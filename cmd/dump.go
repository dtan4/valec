package cmd

import (
	"fmt"
	"strings"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump NAMESPACE",
	Short: "Dump secrets in dotenv format",
	RunE:  doDump,
}

func doDump(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Please specify namespace.")
	}
	namespace := args[0]

	secrets, err := aws.DynamoDB.ListSecrets(tableName, namespace)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve secrets.")
	}

	if len(secrets) == 0 {
		return errors.Errorf("Namespace %s does not exist.", namespace)
	}

	var dotenv []string

	if dotenvTemplate == "" {
		dotenv, err = dumpAll(secrets, quote)
		if err != nil {
			return errors.Wrap(err, "Failed to dump all secrets.")
		}
	} else {
		dotenv, err = dumpWithTemplate(secrets, quote)
		if err != nil {
			return errors.Wrap(err, "Failed to dump secrets with dotenv template.")
		}
	}

	if output == "" {
		for _, line := range dotenv {
			fmt.Println(line)
		}
	} else {
		body := []byte(strings.Join(dotenv, "\n") + "\n")
		if override {
			if err := util.WriteFile(output, body); err != nil {
				return errors.Wrapf(err, "Failed to write dotenv file. filename=%s", output)
			}
		} else {
			if err := util.WriteFileWithoutSection(output, body); err != nil {
				return errors.Wrapf(err, "Failed to write dotenv file. filename=%s", output)
			}
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(dumpCmd)

	dumpCmd.Flags().BoolVar(&override, "override", false, "Override values in existing template")
	dumpCmd.Flags().StringVarP(&output, "output", "o", "", "File to flush dotenv")
	dumpCmd.Flags().BoolVarP(&quote, "quote", "q", false, "Quote values")
	dumpCmd.Flags().StringVarP(&dotenvTemplate, "template", "t", "", "Dotenv template")
}
