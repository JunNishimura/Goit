/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	renameOption string = ""
	deleteOption string = ""
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "handle with branch operation",
	Long:  "handle with branch operation",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		isList, err := cmd.Flags().GetBool("list")
		if err != nil {
			return fmt.Errorf("fail to get list flag: %w", err)
		}
		fmt.Println(isList)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)

	branchCmd.Flags().BoolP("list", "l", false, "show list of branches")
	branchCmd.Flags().StringVarP(&renameOption, "rename", "r", "", "rename branch")
	branchCmd.Flags().StringVarP(&deleteOption, "delete", "d", "", "delete branch")
}
