/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// hashObjectCmd represents the hashObject command
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "calculate the hash of the file",
	Long:  "calculate the hash of the file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hashObject called")
	},
}

func init() {
	rootCmd.AddCommand(hashObjectCmd)
}
