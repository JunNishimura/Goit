/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rFlag bool
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove file from the working tree and the index",
	Long:  "remove file from the working tree and the index",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// args validation
		for _, arg := range args {
			// check if the arg is registered in the Index
			cleanedArg := filepath.Clean(arg)
			cleanedArg = strings.ReplaceAll(cleanedArg, `\`, "/")
			client.Idx.GetEntry([]byte())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	rmCmd.Flags().BoolVarP(&rFlag, "rec", "r", false, "allow recursive removal")
}
