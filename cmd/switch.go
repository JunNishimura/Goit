/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	createOption string
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "switch branches",
	Long:  "switch branches",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// validation
		if len(args) >= 2 {
			return errors.New("fatal: only one reference expected")
		}
		if createOption == "" && len(args) == 0 {
			return errors.New("fatal: missing branch")
		} else if createOption != "" && len(args) >= 1 {
			return errors.New("invalid create option format")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)

	switchCmd.Flags().StringVarP(&createOption, "create", "c", "", "create new branch")
}
