package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/lib"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	yamlExtRegexp = regexp.MustCompile(`\.[yY][aA]?[mM][lL]$`)
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync CONFIGDIR [NAMESPACE]",
	Short: "Synchronize secrets between local file and DynamoDB",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Please specify config file.")
		}
		dirname := args[0]

		files, err := ioutil.ReadDir(dirname)
		if err != nil {
			return errors.Wrapf(err, "Failed to read directory. dirname=%s", dirname)
		}

		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") || !yamlExtRegexp.Match([]byte(file.Name())) {
				continue
			}

			filename := filepath.Join(dirname, file.Name())

			if err := syncFile(filename); err != nil {
				return errors.Wrapf(err, "Failed to synchronize configs. filename=%s", filename)
			}
		}

		return nil
	},
}

func syncFile(filename string) error {
	namespace := yamlExtRegexp.ReplaceAllString(filepath.Base(filename), "")
	fmt.Println(namespace)

	srcConfigs, err := lib.LoadConfigYAML(filename)
	if err != nil {
		return errors.Wrapf(err, "Failed to load configs. filename=%s", filename)
	}

	dstConfigs, err := aws.DynamoDB().ListConfigs(tableName, namespace)
	if err != nil {
		return errors.Wrapf(err, "Failed to retrieve configs. namespace=%s", namespace)
	}

	added, deleted := lib.CompareConfigList(srcConfigs, dstConfigs)
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)

	if len(deleted) > 0 {
		fmt.Printf("%  d configs of %s namespace will be deleted.\n", len(deleted), namespace)
		for _, config := range deleted {
			if noColor {
				fmt.Printf("    - %s\n", config.Key)
			} else {
				red.Printf("    - %s\n", config.Key)
			}
		}

		if !dryRun {
			if err := aws.DynamoDB().Delete(tableName, namespace, deleted); err != nil {
				return errors.Wrapf(err, "Failed to delete configs. namespace=%s", namespace)
			}

			fmt.Printf("  %d configs of %s namespace were successfully deleted.\n", len(deleted), namespace)
		}
	} else {
		fmt.Println("  No config will be deleted.")
	}

	if len(added) > 0 {
		fmt.Printf("  %d configs of %s namespace will be added.\n", len(added), namespace)
		for _, config := range added {
			if noColor {
				fmt.Printf("    + %s\n", config.Key)
			} else {
				green.Printf("    + %s\n", config.Key)
			}
		}

		if !dryRun {
			if err := aws.DynamoDB().Insert(tableName, namespace, added); err != nil {
				return errors.Wrapf(err, "Failed to insert configs. namespace=%s", namespace)
			}

			fmt.Printf("  %d configs of %s namespace were successfully added.\n", len(added), namespace)
		}
	} else {
		fmt.Println("  No config will be added.")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(syncCmd)

	syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Dry run")
}
