/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/JunNishimura/Goit/index"
	"github.com/JunNishimura/Goit/store"
	"github.com/spf13/cobra"
)

var (
	conf        *store.Config
	indexClient *index.Index
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
	// config setting
	cnf, err := store.NewConfig()
	if err != nil {
		fmt.Printf("fail to config setting: %v", err)
		return
	}
	conf = cnf

	// index setting
	index, err := index.NewIndex()
	if err != nil {
		fmt.Printf("fail to NewIndex: %v", err)
		return
	}
	indexClient = index

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func IsGoitInitialized() bool {
	_, err := os.Stat(".goit")
	return !os.IsNotExist(err)
}
