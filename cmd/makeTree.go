/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// makeTreeCmd represents the makeTree command
var makeTreeCmd = &cobra.Command{
	Use:   "make-tree",
	Short: "make tree object from index",
	Long:  "this is a command to make tree object from index",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(makeTreeCmd)
}
