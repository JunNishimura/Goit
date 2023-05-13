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

		// make data from index
		var data []byte
		for _, entry := range indexClient.Entries {
			data = append(data, entry.Path...)
			data = append(data, 0x00)
			data = append(data, entry.Hash...)
		}

		// make tree object
		treeObject := object.NewObject(object.TreeObject, data)

		if err := treeObject.Write(); err != nil {
			return fmt.Errorf("fail to write tree object: %v", err)
		}

		// print out tree object hash
		fmt.Printf("%s\n", treeObject.Hash)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
