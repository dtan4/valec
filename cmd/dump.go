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

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump NAMESPACE",
	Short: "Dump secrets in dotenv format",
	RunE:  doDump,
}

var dumpOpts = struct {
	dotenvTemplate string
	override       bool
	output         string
	quote          bool
}{}

func doDump(cmd *cobra.Command, args []string) error {
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

	var dotenv []string

	if dumpOpts.dotenvTemplate == "" {
		dotenv, err = dumpAll(secrets, dumpOpts.quote)
		if err != nil {
			return errors.Wrap(err, "Failed to dump all secrets")
		}
	} else {
		dotenv, err = dumpWithTemplate(secrets, dumpOpts.quote, dumpOpts.dotenvTemplate, dumpOpts.override)
		if err != nil {
			return errors.Wrap(err, "Failed to dump secrets with dotenv template")
		}
	}

	if dumpOpts.output == "" {
		for _, line := range dotenv {
			fmt.Println(line)
		}
	} else {
		body := []byte(strings.Join(dotenv, "\n") + "\n")

		if _, err := os.Stat(dumpOpts.output); err != nil {
			if os.IsNotExist(err) {
				if err2 := util.WriteFile(dumpOpts.output, body); err != nil {
					return errors.Wrapf(err2, "Failed to write dotenv file %s", dumpOpts.output)
				}
			} else {
				return errors.Wrapf(err, "Failed to open dotenv file %s", dumpOpts.output)
			}
		} else {
			if dumpOpts.override {
				if err := util.WriteFile(dumpOpts.output, body); err != nil {
					return errors.Wrapf(err, "Failed to write dotenv file %s", dumpOpts.output)
				}
			} else {
				if err := util.WriteFileWithoutSection(dumpOpts.output, body); err != nil {
					return errors.Wrapf(err, "Failed to write dotenv file %s", dumpOpts.output)
				}
			}
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(dumpCmd)

	dumpCmd.Flags().BoolVar(&dumpOpts.override, "override", false, "Override values in existing template")
	dumpCmd.Flags().StringVarP(&dumpOpts.output, "output", "o", "", "File to flush dotenv")
	dumpCmd.Flags().BoolVarP(&dumpOpts.quote, "quote", "q", false, "Quote values")
	dumpCmd.Flags().StringVarP(&dumpOpts.dotenvTemplate, "template", "t", "", "Dotenv template")
}
