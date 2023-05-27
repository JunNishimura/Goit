/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/JunNishimura/Goit/internal/object"
	"github.com/JunNishimura/Goit/internal/sha"
	"github.com/spf13/cobra"
)

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "cat goit object",
	Long:  "this is a command to show the goit object",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// args validation check
		if len(args) == 0 {
			return ErrNotSpecifiedHash
		}
		if len(args) > 1 {
			return ErrTooManyArgs
		}

		// get flags
		typeFlag, err := cmd.Flags().GetBool("type")
		if err != nil {
			return fmt.Errorf("fail to get type flag: %w", err)
		}
		printFlag, err := cmd.Flags().GetBool("print")
		if err != nil {
			return fmt.Errorf("fail to get print flag: %w", err)
		}

		// flag check
		if typeFlag && printFlag {
			return ErrIncompatibleFlag
		}

		// get object from hash
		hash, err := sha.ReadHash(args[0])
		if err != nil {
			return ErrInvalidHash
		}
		obj, err := object.GetObject(client.RootGoitPath, hash)
		if err != nil {
			return ErrInvalidHash
		}

		// print object type
		if typeFlag {
			fmt.Printf("%s\n", obj.Type)
		}

		// print object content
		if printFlag {
			if obj.Type == object.TreeObject {
				// need to print out in the different way since hash is written as hexideciaml in data of tree object
				// convert tree object to tree and print out
				tree, err := object.NewTree(client.RootGoitPath, obj)
				if err != nil {
					return fmt.Errorf("fail to get new tree: %w", err)
				}
				fmt.Printf("%s\n", tree)
			} else {
				// hash is written as string
				fmt.Printf("%s\n", obj.Data)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)

	catFileCmd.Flags().BoolP("type", "t", false, "print object type")
	catFileCmd.Flags().BoolP("print", "p", false, "print object content")
}
