/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "commit",
	Long:  "this is a command to commit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("commit called")
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
