package cmd

import (
	"fmt"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/msg"
	"github.com/dtan4/valec/secret"
	"github.com/dtan4/valec/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync SECRETDIR [NAMESPACE]",
	Short: "Synchronize secrets between local file and DynamoDB",
	RunE:  doSync,
}

var syncOpts = struct {
	dryRun bool
}{}

func doSync(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Please specify secret directory.")
	}
	dirname := args[0]

	if rootOpts.noColor {
		msg.DisableColor()
	}

	files, err := util.ListYAMLFiles(dirname)
	if err != nil {
		return errors.Wrapf(err, "Failed to read directory. dirname=%s", dirname)
	}

	srcNamespaces, err := aws.DynamoDB.ListNamespaces(rootOpts.tableName)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve namespaces.")
	}

	dstNamespaces := []string{}

	for _, file := range files {
		namespace, err := util.NamespaceFromPath(file, dirname)
		if err != nil {
			return errors.Wrap(err, "Failed to get namespace.")
		}
		dstNamespaces = append(dstNamespaces, namespace)

		if err := syncFile(file, namespace); err != nil {
			return errors.Wrapf(err, "Failed to synchronize file. filename=%s", file)
		}
	}

	added, deleted := util.CompareStrings(srcNamespaces, dstNamespaces)

	for _, namespace := range added {
		msg.GreenBold.Printf("+ %s\n", namespace)
	}

	if len(added) > 0 {
		fmt.Printf("%d namespaces will be added.\n", len(added))
	}

	for _, namespace := range deleted {
		msg.RedBold.Printf("- %s\n", namespace)
	}

	if len(deleted) > 0 {
		fmt.Printf("%d namespaces will be deleted.\n", len(deleted))
	}

	return nil
}

func syncFile(filename, namespace string) error {
	msg.Bold.Println(namespace)

	_, srcSecrets, err := secret.LoadFromYAML(filename)
	if err != nil {
		return errors.Wrapf(err, "Failed to load secrets. filename=%s", filename)
	}

	dstSecrets, err := aws.DynamoDB.ListSecrets(rootOpts.tableName, namespace)
	if err != nil {
		return errors.Wrapf(err, "Failed to retrieve secrets. namespace=%s", namespace)
	}

	added, updated, deleted := srcSecrets.CompareList(dstSecrets)

	if len(deleted) > 0 {
		fmt.Printf("%  d secrets will be deleted.\n", len(deleted))
		for _, secret := range deleted {
			msg.Red.Printf("    - %s\n", secret.Key)
		}

		if !syncOpts.dryRun {
			if err := aws.DynamoDB.Delete(rootOpts.tableName, namespace, deleted); err != nil {
				return errors.Wrapf(err, "Failed to delete secrets. namespace=%s", namespace)
			}

			fmt.Printf("  %d secrets were successfully deleted.\n", len(deleted))
		}
	}

	if len(updated) > 0 {
		fmt.Printf("  %d secrets will be updated.\n", len(updated))
		for _, secret := range updated {
			msg.Yellow.Printf("    + %s\n", secret.Key)
		}

		if !syncOpts.dryRun {
			if err := aws.DynamoDB.Insert(rootOpts.tableName, namespace, updated); err != nil {
				return errors.Wrapf(err, "Failed to insert secrets. namespace=%s")
			}

			fmt.Printf("  %d secrets were successfully updated.\n", len(updated))
		}
	}

	if len(added) > 0 {
		fmt.Printf("  %d secrets will be added.\n", len(added))
		for _, secret := range added {
			msg.Green.Printf("    + %s\n", secret.Key)
		}

		if !syncOpts.dryRun {
			if err := aws.DynamoDB.Insert(rootOpts.tableName, namespace, added); err != nil {
				return errors.Wrapf(err, "Failed to insert secrets. namespace=%s")
			}

			fmt.Printf("  %d secrets were successfully added.\n", len(added))
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(syncCmd)

	syncCmd.Flags().BoolVar(&syncOpts.dryRun, "dry-run", false, "Dry run")
}
