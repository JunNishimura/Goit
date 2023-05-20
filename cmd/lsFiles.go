/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	isShowStaged bool
)

// lsFilesCmd represents the lsFiles command
var lsFilesCmd = &cobra.Command{
	Use:   "ls-files",
	Short: "print out index",
	Long:  "this is a command to print out index",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !isGoitInitialized() {
			return errors.New("fatal: not a goit repository: .goit")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		for _, entry := range client.Idx.Entries {
			if isShowStaged {
				fmt.Printf("%s    %s\n", entry.Hash, entry.Path)
			} else {
				fmt.Printf("%s\n", entry.Path)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(lsFilesCmd)

	lsFilesCmd.Flags().BoolVarP(&isShowStaged, "staged", "s", false, "show staged contents")
}
