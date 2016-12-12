package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/dtan4/valec/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec NAMESPACE COMMAND [ARG ...]",
	Short: "Execute commands using stored secrets",
	Long: `Execute commands using stored secrets

Stored secrets are consumed as environment variables.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Please specify namespace and command.")
		}
		namespace := args[0]

		secrets, err := aws.DynamoDB.ListSecrets(tableName, namespace)
		if err != nil {
			return errors.Wrapf(err, "Failed to load secrets from DynamoDB. namespace=%s", namespace)
		}

		fetchKeys := map[string]bool{}

		if keys != "" {
			for _, s := range strings.Split(keys, ",") {
				fetchKeys[s] = true
			}
		}

		envs := os.Environ()

		for _, secret := range secrets {
			if len(fetchKeys) > 0 && !fetchKeys[secret.Key] {
				continue
			}

			plainValue, err := aws.KMS.DecryptBase64(secret.Key, secret.Value)
			if err != nil {
				return errors.Wrapf(err, "Failed to decrypt value. key=%q, value=%q", secret.Key, secret.Value)
			}

			envs = append(envs, fmt.Sprintf("%s=%s", secret.Key, plainValue))
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

	execCmd.Flags().StringVarP(&keys, "keys", "k", "", "Secret keys to fetch")
}
