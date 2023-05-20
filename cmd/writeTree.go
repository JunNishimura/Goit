/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/JunNishimura/Goit/object"
	"github.com/spf13/cobra"
)

// writeTreeCmd represents the writeTree command
var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "write tree object from index",
	Long:  "this is a command to write tree object from index",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !IsGoitInitialized() {
			return errors.New("fatal: not a goit repository: .goit")
		}

		// make and write treeObject from index
		rootTreeObject, err := object.WriteTreeObject(client.RootDir, client.Idx.Entries)
		if err != nil {
			return err
		}

		// print out tree object hash
		fmt.Printf("%s\n", rootTreeObject.Hash)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
