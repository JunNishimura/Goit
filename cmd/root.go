/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goit",
	Short: "Git made by Golang",
	Long:  "This is a Git-like CLI tool made by Golang",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func IsGoitInitialized() bool {
	_, err := os.Stat(".goit")
	return !os.IsNotExist(err)
}
