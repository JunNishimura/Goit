/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/JunNishimura/Goit/object"
	"github.com/JunNishimura/Goit/sha"
	"github.com/spf13/cobra"
)

var (
	ErrIncompatibleFlag = errors.New("error: incompatible pair of flags")
	ErrNotSpecifiedHash = errors.New("error: no specified object hash")
	ErrTooManyArgs      = errors.New("error: to many arguments")
	ErrInvalidHash      = errors.New("error: not a valid object hash")
)

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "cat goit object",
	Long:  "this is a command to show the goit object",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !IsGoitInitialized() {
			return errors.New("fatal: not a goit repository: .goit")
		}

		// get flags
		typeFlag, err := cmd.Flags().GetBool("type")
		if err != nil {
			return fmt.Errorf("fail to get type flag: %v", err)
		}
		printFlag, err := cmd.Flags().GetBool("print")
		if err != nil {
			return fmt.Errorf("fail to get print flag: %v", err)
		}

		// flag check
		if typeFlag && printFlag {
			return ErrIncompatibleFlag
		}

		// arguments validation
		if len(args) == 0 {
			return ErrNotSpecifiedHash
		}
		if len(args) > 1 {
			return ErrTooManyArgs
		}

		// get object from hash
		hash, err := sha.ReadHash(args[0])
		if err != nil {
			return ErrInvalidHash
		}
		obj, err := object.GetObject(hash)
		if err != nil {
			return ErrInvalidHash
		}
		fmt.Printf("%v\n", obj)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)

	catFileCmd.Flags().BoolP("type", "t", false, "print object type")
	catFileCmd.Flags().BoolP("print", "p", false, "print object content")
}
