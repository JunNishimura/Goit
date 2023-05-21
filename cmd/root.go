/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/JunNishimura/Goit/store"
	"github.com/spf13/cobra"
)

var (
	client *store.Client
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
	// get client
	tmpClient, err := store.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client = tmpClient

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
