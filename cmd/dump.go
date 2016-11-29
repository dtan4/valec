package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/lib"
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

		if dotenvTemplate == "" {
			if err := dumpAll(configs); err != nil {
				return errors.Wrap(err, "Failed to dump all configs.")
			}
		} else {
			if err := dumpWithTemplate(configs); err != nil {
				return errors.Wrap(err, "Failed to dump configs with dotenv template.")
			}
		}

		return nil
	},
}

func dumpAll(configs []*lib.Config) error {
	for _, config := range configs {
		plainValue, err := aws.KMS().DecryptBase64(config.Key, config.Value)
		if err != nil {
			return errors.Wrap(err, "Failed to decrypt value.")
		}

		fmt.Printf("%s=%s\n", config.Key, plainValue)
	}

	return nil
}

func dumpWithTemplate(configs []*lib.Config) error {
	fp, err := os.Open(dotenvTemplate)
	if err != nil {
		return errors.Wrapf(err, "Failed to open dotenv template. filename=%s", dotenvTemplate)
	}
	defer fp.Close()

	configMap := lib.ConfigsToMap(configs)
	sc := bufio.NewScanner(fp)

	for sc.Scan() {
		line := sc.Text()

		if strings.HasPrefix(line, "#") {
			fmt.Println(line)
			continue
		}

		ss := strings.SplitN(line, "=", 2)
		if len(ss) != 2 {
			fmt.Println(line)
			continue
		}

		key, value := ss[0], ss[1]

		if override || value == "" {
			v, ok := configMap[key]
			if ok {
				plainValue, err := aws.KMS().DecryptBase64(key, v)
				if err != nil {
					return errors.Wrap(err, "Failed to decrypt value.")
				}

				value = plainValue
			}
		}

		fmt.Printf("%s=%s\n", key, value)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(dumpCmd)

	dumpCmd.Flags().BoolVarP(&override, "override", "o", false, "Override values in existing template")
	dumpCmd.Flags().StringVarP(&dotenvTemplate, "template", "t", "", "Dotenv template")
}
