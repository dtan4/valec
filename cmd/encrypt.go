package cmd

import (
	"fmt"
	"strings"

	"github.com/dtan4/valec/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		cipherText, err := aws.KMS().EncryptBase64(keyAlias, key, value)
		if err != nil {
			return errors.Wrapf(err, "Failed to encrypt.")
		}

		fmt.Println(cipherText)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(encryptCmd)
}
