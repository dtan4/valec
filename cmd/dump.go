package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/secret"
	"github.com/dtan4/valec/util"
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

		if len(secrets) == 0 {
			return errors.Errorf("Namespace %s does not exist.", namespace)
		}

		var dotenv []string

		if dotenvTemplate == "" {
			dotenv, err = dumpAll(secrets)
			if err != nil {
				return errors.Wrap(err, "Failed to dump all secrets.")
			}
		} else {
			dotenv, err = dumpWithTemplate(secrets)
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
				util.WriteFile(output, body)
			} else {
				util.WriteFileWithoutSection(output, body)
			}
		}

		return nil
	},
}

func dumpAll(secrets secret.Secrets) ([]string, error) {
	dotenv := []string{}

	for _, secret := range secrets {
		plainValue, err := aws.KMS.DecryptBase64(secret.Key, secret.Value)
		if err != nil {
			return []string{}, errors.Wrap(err, "Failed to decrypt value.")
		}

		dotenv = append(dotenv, fmt.Sprintf("%s=%s", secret.Key, plainValue))
	}

	return dotenv, nil
}

func dumpWithTemplate(secrets secret.Secrets) ([]string, error) {
	fp, err := os.Open(dotenvTemplate)
	if err != nil {
		return []string{}, errors.Wrapf(err, "Failed to open dotenv template. filename=%s", dotenvTemplate)
	}
	defer fp.Close()

	secretMap := secrets.ListToMap()
	sc := bufio.NewScanner(fp)
	dotenv := []string{}

	for sc.Scan() {
		line := sc.Text()

		if strings.HasPrefix(line, "#") {
			dotenv = append(dotenv, line)
			continue
		}

		ss := strings.SplitN(line, "=", 2)
		if len(ss) != 2 {
			dotenv = append(dotenv, line)
			continue
		}

		key, value := ss[0], ss[1]

		if override || value == "" {
			v, ok := secretMap[key]
			if ok {
				plainValue, err := aws.KMS.DecryptBase64(key, v)
				if err != nil {
					return []string{}, errors.Wrap(err, "Failed to decrypt value.")
				}

				value = plainValue
			}
		}

		dotenv = append(dotenv, fmt.Sprintf("%s=%s", key, value))
	}

	return dotenv, nil
}

func init() {
	RootCmd.AddCommand(dumpCmd)

	dumpCmd.Flags().BoolVar(&override, "override", false, "Override values in existing template")
	dumpCmd.Flags().StringVarP(&output, "output", "o", "", "File to flush dotenv")
	dumpCmd.Flags().StringVarP(&dotenvTemplate, "template", "t", "", "Dotenv template")
}
