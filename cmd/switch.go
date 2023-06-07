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
	createOption string
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "switch branches",
	Long:  "switch branches",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// validation
		if len(args) >= 2 {
			return errors.New("fatal: only one reference expected")
		}
		if createOption == "" && len(args) == 0 {
			return errors.New("fatal: missing branch")
		} else if createOption != "" && len(args) >= 1 {
			return errors.New("invalid create option format")
		}

		// switch branch == update HEAD
		if len(args) == 1 {
			if err := client.Head.Update(client.Refs, client.RootGoitPath, args[0]); err != nil {
				return fmt.Errorf("fail to update HEAD: %w", err)
			}
		}

		if createOption != "" {
			if err := client.Refs.AddBranch(client.RootGoitPath, createOption, client.Head.Commit.Hash); err != nil {
				return fmt.Errorf("fail to create new branch %s: %w", createOption, err)
			}
			if err := client.Head.Update(client.Refs, client.RootGoitPath, createOption); err != nil {
				return fmt.Errorf("fail to update HEAD: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)

	switchCmd.Flags().StringVarP(&createOption, "create", "c", "", "create new branch")
}
