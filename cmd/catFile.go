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
	ErrIncompatibleFlag = errors.New("error: incompatible pair of flags")
	ErrNotSpecifiedHash = errors.New("error: no specified object hash")
)

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "cat goit object",
	Long:  "this is a command to show the goit object",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		if len(args) == 0 {
			return ErrNotSpecifiedHash
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)

	catFileCmd.Flags().BoolP("type", "t", false, "print object type")
	catFileCmd.Flags().BoolP("print", "p", false, "print object content")
}
