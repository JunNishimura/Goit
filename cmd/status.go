/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/JunNishimura/Goit/internal/file"
	"github.com/JunNishimura/Goit/internal/object"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "show the working tree status",
	Long:  "show the working tree status",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var statusMessage string

		// set branch info
		statusMessage += fmt.Sprintf("On branch %s\n", client.Head.Reference)

		// walk through working directory
		var newFiles []string
		var modifiedFiles []string
		filePaths, err := file.GetFilePathsUnderDirectoryWithIgnore(".", client.Idx, client.Ignore)
		if err != nil {
			return fmt.Errorf("fail to get files: %w", err)
		}
		for _, filePath := range filePaths {
			_, entry, isRegistered := client.Idx.GetEntry([]byte(filePath))

			if !isRegistered { // new file
				newFiles = append(newFiles, filePath)
			} else {
				// check if the file is modified
				data, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("fail to read %s: %w", filePath, err)
				}
				obj, err := object.NewObject(object.BlobObject, data)
				if err != nil {
					return fmt.Errorf("fail to get new object: %w", err)
				}
				if !entry.Hash.Compare(obj.Hash) {
					modifiedFiles = append(modifiedFiles, filePath)
				}
			}
		}

		// walk through index
		var deletedFiles []string
		for _, entry := range client.Idx.Entries {
			filePath := string(entry.Path)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				deletedFiles = append(deletedFiles, filePath)
			}
		}

		// compare index with HEAD commit
		treeObj, err := object.GetObject(client.RootGoitPath, client.Head.Commit.Tree)
		if err != nil {
			return fmt.Errorf("fail to get tree object: %w", err)
		}
		tree, err := object.NewTree(client.RootGoitPath, treeObj)
		if err != nil {
			return fmt.Errorf("fail to get tree: %w", err)
		}
		diffEntries, err := client.Idx.DiffWithTree(tree)
		if err != nil {
			return fmt.Errorf("fail to get diff entries: %w", err)
		}

		// construct message
		if len(diffEntries) > 0 {
			statusMessage += "\nChanges to be committed:\n  (use 'goit restore --staged <file>...' to unstage)\n"
			for _, diffEntry := range diffEntries {
				statusMessage += color.GreenString("\t%-13s%s\n", diffEntry.Dt, diffEntry.Entry.Path)
			}
		}
		if len(modifiedFiles) > 0 {
			statusMessage += "\nChanges not staged for commit:\n  (use 'goit add/rm <file>...' to update what will be committed)\n  (use 'goit restore <file>...' to discard changes in working directory)\n"
			for _, file := range modifiedFiles {
				statusMessage += color.RedString("\t%-13s%s\n", "modified:", file)
			}
		}
		if len(deletedFiles) > 0 {
			if len(modifiedFiles) == 0 {
				statusMessage += "\nChanges not staged for commit:\n  (use 'goit add/rm <file>...' to update what will be committed)\n  (use 'goit restore <file>...' to discard changes in working directory)\n"
			}
			for _, file := range deletedFiles {
				statusMessage += color.RedString("\t%-13s%s\n", "deleted:", file)
			}
		}
		if len(newFiles) > 0 {
			statusMessage += "\nUntracked files:\n  (use 'goit add <file>...' to include in what will be committed)\n"
			for _, file := range newFiles {
				statusMessage += color.RedString("\t%s\n", file)
			}
		}

		// show message
		fmt.Println(statusMessage)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
