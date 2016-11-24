package cmd

import (
	"fmt"

	"github.com/dtan4/valec/aws"
	"github.com/spf13/cobra"
)

var (
	key  string
	text string
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		plaintText, err := aws.KMS().DecryptBase64(key, text)
		if err != nil {
			return err
		}

		fmt.Println(plaintText)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&key, "key", "k", "", "Key")
	listCmd.Flags().StringVarP(&text, "text", "t", "", "Text to decrypt")
}
