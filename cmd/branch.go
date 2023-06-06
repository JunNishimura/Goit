/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/JunNishimura/Goit/internal/store"
	"github.com/spf13/cobra"
)

var (
	renameOption string = ""
	deleteOption string = ""
)

func deleteBranch(client *store.Client, branchName string) error {
	// branch validation
	if branchName == client.Head.Reference {
		return fmt.Errorf("error: cannot delete current branch '%s'", client.Head.Reference)
	}
	isBranchFound := false
	for _, branch := range client.Refs.Heads {
		if branch.Name == branchName {
			isBranchFound = true
		}
	}
	if !isBranchFound {
		return fmt.Errorf("error: branch '%s' not found", branchName)
	}

	// delete branch
	if err := client.Refs.DeleteBranch(client.RootGoitPath, branchName); err != nil {
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

		// parameter validation
		if !((len(args) == 1 && !isList && renameOption == "" && deleteOption == "") ||
			(len(args) == 0 && isList && renameOption == "" && deleteOption == "") ||
			(len(args) == 0 && !isList && renameOption != "" && deleteOption == "") ||
			(len(args) == 0 && !isList && renameOption == "" && deleteOption != "")) {
			return fmt.Errorf("parameters are not valid")
		}

		// add branch
		if len(args) == 1 {
			addBranchName := args[0]
			addBranchHash := client.Head.Commit.Hash

			if err := client.Refs.AddBranch(client.RootGoitPath, addBranchName, addBranchHash); err != nil {
				return fmt.Errorf("fail to add branch '%s': %w", addBranchName, err)
			}
		}

		// list branches
		if isList {
			client.Refs.ListBranches(client.Head.Reference)
		}

		// rename current branch
		if renameOption != "" {
			if err := client.Refs.RenameBranch(client.Head, client.RootGoitPath, renameOption); err != nil {
				return err
			}
		}

		if deleteOption != "" {
			if err := deleteBranch(client, deleteOption); err != nil {
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
