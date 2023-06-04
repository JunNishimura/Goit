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

func restoreIndex(tree *object.Tree, path string) error {
	// get entry
	_, entry, isEntryFound := client.Idx.GetEntry([]byte(path))
	if !isEntryFound {
		return fmt.Errorf("error: pathspec '%s' did not match any file(s) known to goit", path)
	}

	// get node
	node, isNodeFound := object.GetNode(tree.Children, path)

	// restore index
	if isNodeFound { // if node is in the last commit
		// change hash
		isUpdated, err := client.Idx.Update(client.RootGoitPath, node.Hash, []byte(path))
		if err != nil {
			return fmt.Errorf("fail to update index: %w", err)
		}
		if !isUpdated {
			return errors.New("fail to restore index")
		}
	} else { // if node is not in the last commit
		// delete entry
		if err := client.Idx.DeleteEntry(client.RootGoitPath, entry); err != nil {
			return fmt.Errorf("fail to delete entry: %w", err)
		}
	}

	return nil
}

func restoreWorkingDirectory(path string) error {
	_, entry, isEntryFound := client.Idx.GetEntry([]byte(path))
	if !isEntryFound {
		return fmt.Errorf("error: pathspec '%s' did not match any file(s) known to goit", path)
	}

	obj, err := object.GetObject(client.RootGoitPath, entry.Hash)
	if err != nil {
		return fmt.Errorf("fail to get object '%s': %w", path, err)
	}

	// restore file
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("fail to get abs path '%s': %w", path, err)
	}
	f, err := os.Create(absPath)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, absPath)
	}
	defer f.Close()

	if _, err := f.WriteString(string(obj.Data)); err != nil {
		return fmt.Errorf("fail to write to file '%s': %w", absPath, err)
	}

	return nil
}

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
			branchPath := filepath.Join(client.RootGoitPath, "refs", "heads", client.Head.Reference)
			if _, err := os.Stat(branchPath); os.IsNotExist(err) {
				return errors.New("fatal: could not resolve HEAD")
			}

			// get tree from HEAD commit
			treeObject, err := object.GetObject(client.RootGoitPath, client.Head.Commit.Tree)
			if err != nil {
				return fmt.Errorf("fail to get tree object from commit HEAD: %w", err)
			}
			tree, err := object.NewTree(client.RootGoitPath, treeObject)
			if err != nil {
				return fmt.Errorf("fail to get tree: %w", err)
			}

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

						// restore index
						if err := restoreIndex(tree, cleanedRelPath); err != nil {
							return err
						}
					}
				} else { // file
					cleanedArg := filepath.Clean(arg)
					cleanedArg = strings.ReplaceAll(cleanedArg, `\`, "/")

					// restore index
					if err := restoreIndex(tree, cleanedArg); err != nil {
						return err
					}
				}
			}
		} else {
			// execute restore working directory
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

						// restore working directory
						if err := restoreWorkingDirectory(cleanedRelPath); err != nil {
							return err
						}
					}
				} else { // file
					cleanedArg := filepath.Clean(arg)
					cleanedArg = strings.ReplaceAll(cleanedArg, `\`, "/")

					// restore working directory
					if err := restoreWorkingDirectory(cleanedArg); err != nil {
						return err
					}
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().Bool("staged", false, "restore index")
}
