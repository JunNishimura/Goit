/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JunNishimura/Goit/internal/file"
	"github.com/JunNishimura/Goit/internal/object"
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

		// get staged option
		isStaged, err := cmd.Flags().GetBool("staged")
		if err != nil {
			return fmt.Errorf("fail to get staged flag: %w", err)
		}

		// staged validation check
		if isStaged {
			// restore --stage is comparing index with commit object pointed by HEAD
			// so, at lease one commit is needed
			branchPath := filepath.Join(client.RootGoitPath, "refs", "heads", string(client.Head))
			if _, err := os.Stat(branchPath); os.IsNotExist(err) {
				return errors.New("fatal: could not resolve HEAD")
			}

			// check if args are registered in index
			for _, arg := range args {
				argAbsPath, err := filepath.Abs(arg)
				if err != nil {
					return fmt.Errorf("fail to get arg abs path: %w", err)
				}
				f, err := os.Stat(argAbsPath)
				if err != nil {
					return fmt.Errorf("%w: %s", ErrIOHandling, argAbsPath)
				}

				if f.IsDir() { // directory
					filePaths, err := file.GetFilePathsUnderDirectory(argAbsPath)
					if err != nil {
						return fmt.Errorf("fail to get file path under directory: %w", err)
					}
					for _, filePath := range filePaths {
						curPath, err := os.Getwd()
						if err != nil {
							return fmt.Errorf("fail to get current directory: %w", err)
						}
						relPath, err := filepath.Rel(curPath, filePath)
						if err != nil {
							return fmt.Errorf("fail to get relative path: %w", err)
						}
						cleanedRelPath := strings.ReplaceAll(relPath, `\`, "/")
						if !client.Idx.IsPathStaged([]byte(cleanedRelPath)) {
							return fmt.Errorf("error: pathspec '%s' did not match any file(s) known to goit", cleanedRelPath)
						}
					}
				} else { // file
					cleanedArg := filepath.Clean(arg)
					cleanedArg = strings.ReplaceAll(cleanedArg, `\`, "/")
					if !client.Idx.IsPathStaged([]byte(cleanedArg)) {
						return fmt.Errorf("error: pathspec '%s' did not match any file(s) known to goit", cleanedArg)
					}
				}
			}

			// get HEAD commit
			headCommit, err := getHeadCommit()
			if err != nil {
				return fmt.Errorf("fail to get HEAD commit: %w", err)
			}

			// get tree from HEAD commit
			treeObject, err := object.GetObject(client.RootGoitPath, headCommit.Tree)
			if err != nil {
				return fmt.Errorf("fail to get tree object from commit HEAD: %w", err)
			}
			tree, err := object.NewTree(client.RootGoitPath, treeObject)
			if err != nil {
				return fmt.Errorf("fail to get tree: %w", err)
			}
			fmt.Println(tree)

			// search
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().Bool("staged", false, "restore index")
}
