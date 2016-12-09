package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/secret"
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

		secrets, err := aws.DynamoDB.ListSecrets(tableName, namespace)
		if err != nil {
			return errors.Wrap(err, "Failed to retrieve secrets.")
		}

		if dotenvTemplate == "" {
			if err := dumpAll(secrets); err != nil {
				return errors.Wrap(err, "Failed to dump all secrets.")
			}
		} else {
			if err := dumpWithTemplate(secrets); err != nil {
				return errors.Wrap(err, "Failed to dump secrets with dotenv template.")
			}
		}

		return nil
	},
}

func dumpAll(secrets []*secret.Secret) error {
	for _, secret := range secrets {
		plainValue, err := aws.KMS.DecryptBase64(secret.Key, secret.Value)
		if err != nil {
			return errors.Wrap(err, "Failed to decrypt value.")
		}

		fmt.Printf("%s=%s\n", secret.Key, plainValue)
	}

	return nil
}

func dumpWithTemplate(secrets []*secret.Secret) error {
	fp, err := os.Open(dotenvTemplate)
	if err != nil {
		return errors.Wrapf(err, "Failed to open dotenv template. filename=%s", dotenvTemplate)
	}
	defer fp.Close()

	secretMap := secret.ListToMap(secrets)
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
			v, ok := secretMap[key]
			if ok {
				plainValue, err := aws.KMS.DecryptBase64(key, v)
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
