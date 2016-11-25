package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/dtan4/valec/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Please specify namespace or config file (-f FILE).")
		}
		namespace := args[0]

		configs, err := aws.DynamoDB().List(tableName, namespace)
		if err != nil {
			return errors.Wrapf(err, "Failed to load configs from DynamoDB. namespace=%s", namespace)
		}

		envs := os.Environ()

		for _, config := range configs {
			plainValue, err := aws.KMS().DecryptBase64(config.Key, config.Value)
			if err != nil {
				return errors.Wrapf(err, "Failed to decrypt value. key=%q, value=%q", config.Key, config.Value)
			}

			envs = append(envs, fmt.Sprintf("%s=%s", config.Key, plainValue))
		}

		execCmd := exec.Command(args[1], args[2:]...)
		execCmd.Env = envs
		execCmd.Stderr = os.Stderr
		execCmd.Stdout = os.Stdout
		err = execCmd.Run()

		if execCmd.Process == nil {
			return errors.Wrap(err, "Failed to execute command.")
		}

		os.Exit(execCmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(execCmd)
}
