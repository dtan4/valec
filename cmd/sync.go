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

		srcConfigs, err := lib.LoadConfigYAML(filename)
		if err != nil {
			return errors.Wrapf(err, "Failed to load configs. filename=%s", filename)
		}

		var namespace string

		if len(args) == 1 {
			namespace = yamlExtRegexp.ReplaceAllString(filepath.Base(filename), "")
		} else {
			namespace = args[1]
		}

		dstConfigs, err := aws.DynamoDB().ListConfigs(tableName, namespace)
		if err != nil {
			return errors.Wrapf(err, "Failed to retrieve configs. namespace=%s", namespace)
		}

		added, deleted := lib.CompareConfigList(srcConfigs, dstConfigs)

		if err := aws.DynamoDB().Delete(tableName, namespace, deleted); err != nil {
			return errors.Wrapf(err, "Failed to delete configs. namespace=%s", namespace)
		}

		if len(deleted) > 0 {
			fmt.Printf("%d configs of %q namespace were successfully deleted!\n", len(deleted), namespace)
		} else {
			fmt.Println("No config was deleted.")
		}

		if err := aws.DynamoDB().Insert(tableName, namespace, added); err != nil {
			return errors.Wrapf(err, "Failed to insert configs. namespace=%s", namespace)
		}

		if len(added) > 0 {
			fmt.Printf("%d configs of %q namespace were successfully added!\n", len(added), namespace)
		} else {
			fmt.Println("No config was added.")
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
