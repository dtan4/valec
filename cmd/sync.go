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

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync CONFIGFILE [NAMESPACE]",
	Short: "Synchronize secrets between local file and DynamoDB",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Please specify config file.")
		}
		filename := args[0]

		configs, err := lib.LoadConfigYAML(filename)
		if err != nil {
			return errors.Wrapf(err, "Failed to load configs. filename=%s", filename)
		}

		var namespace string

		if len(args) == 1 {
			namespace = yamlExtRegexp.ReplaceAllString(filepath.Base(filename), "")
		} else {
			namespace = args[1]
		}

		if err := aws.DynamoDB().Insert(tableName, namespace, configs); err != nil {
			return errors.Wrapf(err, "Failed to insert configs. namespace=%s", namespace)
		}

		fmt.Printf("%d configs of %q namespace are successfully synchronized!\n", len(configs), namespace)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
