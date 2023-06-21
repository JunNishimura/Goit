/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
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
			if _, _, isRegistered := client.Idx.GetEntry([]byte(cleanedArg)); !isRegistered {
				return fmt.Errorf("fatal: pathspec '%s' did not match any files", arg)
			}
		}

		// remove file from working tree and index
		for _, arg := range args {
			cleanedArg := filepath.Clean(arg)
			cleanedArg = strings.ReplaceAll(cleanedArg, `\`, "/")

			// remove from the working tree
			if _, err := os.Stat(cleanedArg); !os.IsNotExist(err) {
				if err := os.Remove(cleanedArg); err != nil {
					return fmt.Errorf("fail to delete %s from the working tree: %w", cleanedArg, err)
				}
			}

			// remove from the index
			if err := client.Idx.DeleteEntry(client.RootGoitPath, []byte(cleanedArg)); err != nil {
				return fmt.Errorf("fail to delete '%s' from the index: %w", cleanedArg, err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	rmCmd.Flags().BoolVarP(&rFlag, "rec", "r", false, "allow recursive removal")
}
