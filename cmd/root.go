package cmd

import (
	"fmt"
	"os"

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
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if debug {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Println(err)
		}
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug mode")
	RootCmd.PersistentFlags().StringVar(&keyAlias, "key", defaultKeyAlias, "KMS key alias")
	RootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colorize output")
	RootCmd.PersistentFlags().StringVar(&tableName, "table-name", defaultTableName, "DynamoDB table name")
	RootCmd.PersistentFlags().StringVar(&region, "region", "", "AWS region")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
