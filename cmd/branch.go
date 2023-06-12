/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/JunNishimura/Goit/internal/log"
	"github.com/spf13/cobra"
)

var (
	renameOption string = ""
	deleteOption string = ""
)

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
			if err := gLogger.WriteBranch(log.NewRecord(log.BranchRecord, nil, addBranchHash, client.Conf.GetUserName(), client.Conf.GetEmail(), time.Now(), fmt.Sprintf("Created from %s", client.Head.Reference)), addBranchName); err != nil {
				return fmt.Errorf("log error: %w", err)
			}
		}

		// list branches
		if isList {
			client.Refs.ListBranches(client.Head.Reference)
		}

		// rename current branch
		if renameOption != "" {
			if err := client.Refs.RenameBranch(client.RootGoitPath, client.Head.Reference, renameOption); err != nil {
				return fmt.Errorf("fail to rename branch: %w", err)
			}
			// update HEAD
			if err := client.Head.Update(client.Refs, client.RootGoitPath, renameOption); err != nil {
				return fmt.Errorf("fail to update HEAD: %w", err)
			}
		}

		if deleteOption != "" {
			if err := client.Refs.DeleteBranch(client.RootGoitPath, client.Head.Reference, deleteOption); err != nil {
				return fmt.Errorf("fail to delete branch: %w", err)
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
