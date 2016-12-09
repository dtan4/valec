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
	Use:   "encrypt KEY=VALUE",
	Short: "Encrypt secret",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Please specify KEY=VALUE.")
		}
		keyValue := args[0]

		ss := strings.SplitN(keyValue, "=", 2)
		if len(ss) < 2 {
			return errors.Errorf("Given argument is invalid format, should be KEY=VALUE. argument=%q", keyValue)
		}
		key, value := ss[0], ss[1]

		cipherText, err := aws.KMS.EncryptBase64(keyAlias, key, value)
		if err != nil {
			return errors.Wrapf(err, "Failed to encrypt.")
		}

		if secretFile == "" {
			fmt.Println(cipherText)
		} else {
			secretMap := map[string]string{}

			if _, err := os.Stat(secretFile); err == nil {
				secrets, err2 := secret.LoadFromYAML(secretFile)
				if err2 != nil {
					return errors.Wrapf(err2, "Failed to load local secret file. filename=%s", secretFile)
				}

				secretMap = secrets.ListToMap()
			}

			secretMap[key] = cipherText
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
