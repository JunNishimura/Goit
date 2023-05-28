/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "restore file",
	Long:  "restore file",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// args validation
		if len(args) == 0 {
			return errors.New("fatal: you must specify path(s) to restore")
		}

		// check if arg exists
		for _, arg := range args {
			argPath, err := filepath.Abs(arg)
			if err != nil {
				return fmt.Errorf("fail to get absolute path of %s: %w", arg, err)
			}
			if _, err := os.Stat(argPath); os.IsNotExist(err) {
				return fmt.Errorf("error: pathspec '%s' did not match any file(s) known to goit", arg)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().Bool("staged", false, "restore index")
}
