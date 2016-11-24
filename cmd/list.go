package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/lib"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
		if len(args) != 1 {
			return fmt.Errorf("Please specify config file.")
		}
		filename := args[0]

		var configs []*lib.Config

		body, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(body, &configs); err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(w, strings.Join([]string{"KEY", "VALUE"}, "\t"))

		for _, config := range configs {
			plainValue, err := aws.KMS().DecryptBase64(config.Key, config.Value)
			if err != nil {
				return err
			}

			fmt.Fprintln(w, strings.Join([]string{config.Key, plainValue}, "\t"))
		}

		w.Flush()

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
