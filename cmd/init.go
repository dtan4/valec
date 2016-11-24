// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//

package cmd

import (
	"fmt"

	"github.com/dtan4/valec/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := aws.DynamoDB().CreateTable(tableName); err != nil {
			return errors.Wrapf(err, "Failed to create DynamoDB table. table=%s", tableName)
		}

		fmt.Printf("DynamoDB table %s successfully created!\n", tableName)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
