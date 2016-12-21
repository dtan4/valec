package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/secret"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt KEY1=VALUE1 [KEY2=VALUE2 ...]",
	Short: "Encrypt secret",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Please specify KEY=VALUE.")
		}

		newSecretMap := map[string]string{}

		if args[0] == "-" {
			lines := scanFromStdin(os.Stdin)
			for _, line := range lines {
				ss := strings.SplitN(line, "=", 2)
				if len(ss) < 2 {
					continue
				}
				key, value := ss[0], ss[1]

				cipherText, err := aws.KMS.EncryptBase64(keyAlias, key, value)
				if err != nil {
					return errors.Wrapf(err, "Failed to encrypt.")
				}

				newSecretMap[key] = cipherText
			}
		} else {
			for _, arg := range args {
				ss := strings.SplitN(arg, "=", 2)
				if len(ss) < 2 {
					return errors.Errorf("Given argument is invalid format, should be KEY=VALUE. argument=%q", arg)
				}
				key, value := ss[0], ss[1]

				cipherText, err := aws.KMS.EncryptBase64(keyAlias, key, value)
				if err != nil {
					return errors.Wrapf(err, "Failed to encrypt.")
				}

				newSecretMap[key] = cipherText
			}
		}

		if secretFile == "" {
			for _, v := range newSecretMap {
				fmt.Println(v)
			}
		} else {
			secretMap := map[string]string{}

			if _, err := os.Stat(secretFile); err == nil {
				secrets, err2 := secret.LoadFromYAML(secretFile)
				if err2 != nil {
					return errors.Wrapf(err2, "Failed to load local secret file. filename=%s", secretFile)
				}

				secretMap = secrets.ListToMap()
			}

			for k, v := range newSecretMap {
				secretMap[k] = v
			}
			newSecrets := secret.MapToList(secretMap)

			if err := newSecrets.SaveAsYAML(secretFile); err != nil {
				return errors.Wrapf(err, "Failed to update local secret file. filename=%s", secretFile)
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringVar(&secretFile, "add", "", "Add to local secret file")
}
