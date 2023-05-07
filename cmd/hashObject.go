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

// hashObjectCmd represents the hashObject command
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "calculate the hash of the file",
	Long:  "calculate the hash of the file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !IsGoitInitialized() {
			return errors.New("fatal: not a goit repository: .goit")
		}

		for _, arg := range args {
			obj, err := object.NewBlobObject(arg)
			if err != nil {
				return err
			}

			fmt.Println(obj.Hash.String())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(hashObjectCmd)
}
