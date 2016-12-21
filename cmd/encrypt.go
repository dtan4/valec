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
	Use:   "encrypt [KEY1=VALUE1 [KEY2=VALUE2 ...]] [-]",
	Short: "Encrypt secret",
	Long: `Encrypt secret

Read from command line arguments:
  $ valec encrypt KEY1=VALUE1 KEY2=VALUE2

Read from stdin:
  $ cat .env
  KEY1=VALUE1
  KEY2=VALUE2
  $ cat .env | valec encrypt -`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Please specify KEY=VALUE.")
		}

		newSecretMap := map[string]string{}
		var err error

		if args[0] == "-" {
			newSecretMap, err = readFromStdin()
			if err != nil {
				return errors.Wrap(err, "Failed to read secret from stdin.")
			}
		} else {
			if interactive {
				newSecretMap, err = readFromArgsInteractive(args)
				if err != nil {
					return errors.Wrap(err, "Failed to read secret from args.")
				}
			} else {
				newSecretMap, err = readFromArgs(args)
				if err != nil {
					return errors.Wrap(err, "Failed to read secret from args.")
				}
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

func readFromStdin() (map[string]string, error) {
	secretMap := map[string]string{}
	lines := scanLines(os.Stdin)

	for _, line := range lines {
		ss := strings.SplitN(line, "=", 2)
		if len(ss) < 2 {
			continue
		}
		key, value := ss[0], ss[1]

		cipherText, err := aws.KMS.EncryptBase64(keyAlias, key, value)
		if err != nil {
			return map[string]string{}, errors.Wrapf(err, "Failed to encrypt secret. key=%s", key)
		}

		secretMap[key] = cipherText
	}

	return secretMap, nil
}

func readFromArgs(args []string) (map[string]string, error) {
	secretMap := map[string]string{}

	for _, arg := range args {
		ss := strings.SplitN(arg, "=", 2)
		if len(ss) < 2 {
			return map[string]string{}, errors.Errorf("Given argument is invalid format, should be KEY=VALUE. argument=%q", arg)
		}
		key, value := ss[0], ss[1]

		cipherText, err := aws.KMS.EncryptBase64(keyAlias, key, value)
		if err != nil {
			return map[string]string{}, errors.Wrapf(err, "Failed to encrypt secret. key=%s", key)
		}

		secretMap[key] = cipherText
	}

	return secretMap, nil
}

func readFromArgsInteractive(args []string) (map[string]string, error) {
	secretMap := map[string]string{}

	for _, arg := range args {
		key := arg
		fmt.Printf("%s: ", arg)

		value, err := scanNoEcho()
		if err != nil {
			return map[string]string{}, errors.Wrap(err, "Failed to read secret value.")
		}

		cipherText, err := aws.KMS.EncryptBase64(keyAlias, key, value)
		if err != nil {
			return map[string]string{}, errors.Wrapf(err, "Failed to encrypt secret. key=%s", key)
		}

		secretMap[key] = cipherText
	}

	return secretMap, nil
}

func init() {
	RootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringVar(&secretFile, "add", "", "Add to local secret file")
	encryptCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive value input")
}
