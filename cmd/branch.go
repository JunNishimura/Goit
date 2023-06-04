/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JunNishimura/Goit/internal/store"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	renameOption string = ""
	deleteOption string = ""
)

func list(client *store.Client) {
	for _, branch := range client.Refs.Heads {
		if branch.Name == client.Head.Reference {
			color.Green("* %s", branch.Name)
		} else {
			fmt.Println(branch.Name)
		}
	}
}

func rename(client *store.Client, newName string) error {
	// check if new name is not used for other branches
	for _, branch := range client.Refs.Heads {
		if branch.Name == newName {
			return fmt.Errorf("fatal: branch named '%s' already exists", newName)
		}
	}

	// rename current branch
	branch, err := client.Refs.GetBranch(client.Head.Reference)
	if err != nil {
		return err
	}
	branch.Name = newName

	// rename file
	oldPath := filepath.Join(client.RootGoitPath, "refs", "heads", client.Head.Reference)
	newPath := filepath.Join(client.RootGoitPath, "refs", "heads", newName)
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("fail to rename: %w", err)
	}

	// rename HEAD
	if err := client.Head.Update(client.RootGoitPath, newName); err != nil {
		return err
	}

	return nil
}

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "handle with branch operation",
	Long:  "handle with branch operation",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// get list flag
		isList, err := cmd.Flags().GetBool("list")
		if err != nil {
			return fmt.Errorf("fail to get list flag: %w", err)
		}

		// flag validation
		if !((isList && renameOption == "" && deleteOption == "") || (!isList && renameOption != "" && deleteOption == "") || (!isList && renameOption == "" && deleteOption != "")) {
			return fmt.Errorf("only one flag is acceptible")
		}

		// list branches
		if isList {
			list(client)
		}

		// rename current branch
		if renameOption != "" {
			if err := rename(client, renameOption); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)

	branchCmd.Flags().BoolP("list", "l", false, "show list of branches")
	branchCmd.Flags().StringVarP(&renameOption, "rename", "r", "", "rename branch")
	branchCmd.Flags().StringVarP(&deleteOption, "delete", "d", "", "delete branch")
}
