package cmd

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/lib"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	yamlExtRegexp = regexp.MustCompile(`\.[yY][aA]?[mM][lL]$`)
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save secrets in local file to DynamoDB",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Please specify config file.")
		}
		filename := args[0]

		configs, err := lib.LoadConfigYAML(filename)
		if err != nil {
			return errors.Wrapf(err, "Failed to load configs. filename=%s", filename)
		}

		namespace := yamlExtRegexp.ReplaceAllString(filepath.Base(filename), "")

		if err := aws.DynamoDB().Insert(tableName, namespace, configs); err != nil {
			return errors.Wrapf(err, "Failed to insert configs. namespace=%s", namespace)
		}

		fmt.Printf("%d configs of %q namespace are successfully saved!\n", len(configs), namespace)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(saveCmd)
}
