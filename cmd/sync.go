package cmd

import (
	"fmt"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/secret"
	"github.com/dtan4/valec/util"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync SECRETDIR [NAMESPACE]",
	Short: "Synchronize secrets between local file and DynamoDB",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Please specify secret directory.")
		}
		dirname := args[0]

		files, err := util.ListYAMLFiles(dirname)
		if err != nil {
			return errors.Wrapf(err, "Failed to read directory. dirname=%s", dirname)
		}

		for _, file := range files {
			if err := syncFile(file, dirname); err != nil {
				return errors.Wrapf(err, "Failed to synchronize file. filename=%s", file)
			}
		}

		return nil
	},
}

func syncFile(filename, dirname string) error {
	namespace := util.NamespaceFromPath(filename, dirname)

	if noColor {
		fmt.Println(namespace)
	} else {
		color.New(color.Bold).Println(namespace)
	}

	_, srcSecrets, err := secret.LoadFromYAML(filename)
	if err != nil {
		return errors.Wrapf(err, "Failed to load secrets. filename=%s", filename)
	}

	dstSecrets, err := aws.DynamoDB.ListSecrets(tableName, namespace)
	if err != nil {
		return errors.Wrapf(err, "Failed to retrieve secrets. namespace=%s", namespace)
	}

	added, updated, deleted := srcSecrets.CompareList(dstSecrets)
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	if len(deleted) > 0 {
		fmt.Printf("%  d secrets will be deleted.\n", len(deleted))
		for _, secret := range deleted {
			if noColor {
				fmt.Printf("    - %s\n", secret.Key)
			} else {
				red.Printf("    - %s\n", secret.Key)
			}
		}

		if !dryRun {
			if err := aws.DynamoDB.Delete(tableName, namespace, deleted); err != nil {
				return errors.Wrapf(err, "Failed to delete secrets. namespace=%s", namespace)
			}

			fmt.Printf("  %d secrets were successfully deleted.\n", len(deleted))
		}
	}

	if len(updated) > 0 {
		fmt.Printf("  %d secrets will be updated.\n", len(updated))
		for _, secret := range updated {
			if noColor {
				fmt.Printf("    + %s\n", secret.Key)
			} else {
				yellow.Printf("    + %s\n", secret.Key)
			}
		}

		if !dryRun {
			if err := aws.DynamoDB.Insert(tableName, namespace, updated); err != nil {
				return errors.Wrapf(err, "Failed to insert secrets. namespace=%s")
			}

			fmt.Printf("  %d secrets were successfully updated.\n", len(updated))
		}
	}

	if len(added) > 0 {
		fmt.Printf("  %d secrets will be added.\n", len(added))
		for _, secret := range added {
			if noColor {
				fmt.Printf("    + %s\n", secret.Key)
			} else {
				green.Printf("    + %s\n", secret.Key)
			}
		}

		if !dryRun {
			if err := aws.DynamoDB.Insert(tableName, namespace, added); err != nil {
				return errors.Wrapf(err, "Failed to insert secrets. namespace=%s")
			}

			fmt.Printf("  %d secrets were successfully added.\n", len(added))
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(syncCmd)

	syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Dry run")
}
