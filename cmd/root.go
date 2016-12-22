package cmd

import (
	"fmt"
	"os"

	"github.com/dtan4/valec/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	defaultKeyAlias  = "valec"
	defaultTableName = "valec"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "valec",
	Short:         "Handle application secrets securely",
	Long: `Valec is a CLI tool to handle application secrets securely using AWS DynamoDB and KMS.
Valec enables you to manage application secrets in your favorite VCS.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := aws.Initialize(rootOpts.region); err != nil {
			return errors.Wrap(err, "Failed to initialize AWS API clients.")
		}

		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

var rootOpts = struct {
	debug     bool
	keyAlias  string
	noColor   bool
	tableName string
	region    string
}{}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if rootOpts.debug {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Println(err)
		}
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVar(&rootOpts.debug, "debug", false, "Debug mode")
	RootCmd.PersistentFlags().StringVar(&rootOpts.keyAlias, "key", defaultKeyAlias, "KMS key alias")
	RootCmd.PersistentFlags().BoolVar(&rootOpts.noColor, "no-color", false, "Disable colorized output")
	RootCmd.PersistentFlags().StringVar(&rootOpts.tableName, "table-name", defaultTableName, "DynamoDB table name")
	RootCmd.PersistentFlags().StringVar(&rootOpts.region, "region", "", "AWS region")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
