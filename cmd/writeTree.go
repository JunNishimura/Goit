/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/JunNishimura/Goit/internal/object"
	"github.com/spf13/cobra"
)

// writeTreeCmd represents the writeTree command
var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "write tree object from index",
	Long:  "this is a command to write tree object from index",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// make and write treeObject from index
		rootTreeObject, err := object.WriteTreeObject(client.RootGoitPath, client.Idx.Entries)
		if err != nil {
			return fmt.Errorf("fail to write tree object: %w", err)
		}

		// print out tree object hash
		fmt.Printf("%s\n", rootTreeObject.Hash)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
