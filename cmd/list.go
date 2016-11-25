package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/lib"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var configFile string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [NAMESPACE]",
	Short: "List stored secrets",
	Long: `List stored secrets

To list secrets stored in DynamoDB, specify namespace:
  $ valec list NAMESPACE
To list secrets stored in local file, specify file:
  $ valec list -f qa.yaml

Encrypted values are decrypted and printed as plain text.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			configs []*lib.Config
			err     error
		)

		if configFile == "" {
			if len(args) != 1 {
				return errors.New("Please specify namespace or config file (-f FILE).")
			}
			namespace := args[0]

			configs, err = aws.DynamoDB().List(tableName, namespace)
			if err != nil {
				return errors.Wrapf(err, "Failed to load configs from DynamoDB. namespace=%s", namespace)
			}
		} else {
			configs, err = lib.LoadConfigYAML(configFile)
			if err != nil {
				return errors.Wrapf(err, "Failed to load configs from file. filename=%s", configFile)
			}
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(w, strings.Join([]string{"KEY", "VALUE"}, "\t"))

		for _, config := range configs {
			plainValue, err := aws.KMS().DecryptBase64(config.Key, config.Value)
			if err != nil {
				return errors.Wrapf(err, "Failed to decrypt value. key=%q, value=%q", config.Key, config.Value)
			}

			fmt.Fprintln(w, strings.Join([]string{config.Key, plainValue}, "\t"))
		}

		w.Flush()

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&configFile, "file", "f", "", "Config file")
}
