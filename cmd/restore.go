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
	"github.com/JunNishimura/Goit/internal/store"
	"github.com/spf13/cobra"
)

func restoreIndex(rootGoitPath, path string, index *store.Index, tree *object.Tree) error {
	// get entry
	_, _, isEntryFound := index.GetEntry([]byte(path))

	// get node
	node, isNodeFound := object.GetNode(tree.Children, path)

	// if the path is registered in the Index
	if isEntryFound {
		// restore index
		if isNodeFound { // if the file is updated
			// change hash
			isUpdated, err := index.Update(rootGoitPath, node.Hash, []byte(path))
			if err != nil {
				return fmt.Errorf("fail to update index: %w", err)
			}
			if !isUpdated {
				return errors.New("fail to restore index")
			}
		} else { // if the file is newly added
			// delete entry
			if err := index.DeleteEntry(rootGoitPath, []byte(path)); err != nil {
				return fmt.Errorf("fail to delete entry: %w", err)
			}
		}
	} else { // if the path is not registered in the index,
		if isNodeFound { // if the file is deleted
			isUpdated, err := index.Update(rootGoitPath, node.Hash, []byte(path))
			if err != nil {
				return fmt.Errorf("fail to update index: %w", err)
			}
			if !isUpdated {
				return errors.New("fail to restore index")
			}
		} else {
			return fmt.Errorf("error: pathspec '%s' did not match any file(s) known to goit", path)
		}
	}

	return nil
}

func restoreWorkingDirectory(rootGoitPath, path string, index *store.Index) error {
	_, entry, isEntryFound := index.GetEntry([]byte(path))
	if !isEntryFound {
		return fmt.Errorf("error: pathspec '%s' did not match any file(s) known to goit", path)
	}

	obj, err := object.GetObject(rootGoitPath, entry.Hash)
	if err != nil {
		return fmt.Errorf("fail to get object '%s': %w", path, err)
	}

	// get abs path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("fail to get abs path '%s': %w", path, err)
	}

	// make sure the parent path exists
	parentPath := filepath.Dir(absPath)
	if _, err := os.Stat(parentPath); os.IsNotExist(err) {
		if err := os.MkdirAll(parentPath, 0777); err != nil {
			return fmt.Errorf("fail to make directory %s: %w", parentPath, err)
		}
	}

	// restore file
	f, err := os.Create(absPath)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, absPath)
	}
	defer f.Close()

	// write contents to the file
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
				if os.IsNotExist(err) { // even if the file is not found, the file might be the deleted file
					// get node
					cleanedArg := filepath.Clean(arg)
					cleanedArg = strings.ReplaceAll(cleanedArg, `\`, "/")
					node, isNodeFound := object.GetNode(tree.Children, cleanedArg)
					if !isNodeFound {
						return fmt.Errorf("error: pathspec '%s' did not match any file(s) known to goit", arg)
					}

					// check if the arg is dir or not
					if len(node.Children) > 0 { // node is directory
						paths := node.GetPaths()

						for _, path := range paths {
							if err := restoreIndex(client.RootGoitPath, path, client.Idx, tree); err != nil {
								return err
							}
						}
					} else { // node is a file
						if err := restoreIndex(client.RootGoitPath, cleanedArg, client.Idx, tree); err != nil {
							return err
						}
					}

					continue
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
						if err := restoreIndex(client.RootGoitPath, cleanedRelPath, client.Idx, tree); err != nil {
							return err
						}
					}
				} else { // file
					cleanedArg := filepath.Clean(arg)
					cleanedArg = strings.ReplaceAll(cleanedArg, `\`, "/")

					// restore index
					if err := restoreIndex(client.RootGoitPath, cleanedArg, client.Idx, tree); err != nil {
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
				if os.IsNotExist(err) {
					// check if the arg is registered in the index
					cleanedArg := filepath.Clean(arg)
					cleanedArg = strings.ReplaceAll(cleanedArg, `\`, "/")
					_, _, isRegistered := client.Idx.GetEntry([]byte(cleanedArg))
					isRegisteredAsDir := client.Idx.IsRegisteredAsDirectory(cleanedArg)

					if !(isRegistered || isRegisteredAsDir) {
						return fmt.Errorf("error: pathspec '%s' did not match any file(s) known to goit", arg)
					}

					if isRegisteredAsDir {
						entries := client.Idx.GetEntriesByDirectory(cleanedArg)
						for _, entry := range entries {
							if err := restoreWorkingDirectory(client.RootGoitPath, string(entry.Path), client.Idx); err != nil {
								return err
							}
						}
					} else {
						if err := restoreWorkingDirectory(client.RootGoitPath, cleanedArg, client.Idx); err != nil {
							return err
						}
					}

					continue
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
						if err := restoreWorkingDirectory(client.RootGoitPath, cleanedRelPath, client.Idx); err != nil {
							return err
						}
					}
				} else { // file
					cleanedArg := filepath.Clean(arg)
					cleanedArg = strings.ReplaceAll(cleanedArg, `\`, "/")

					// restore working directory
					if err := restoreWorkingDirectory(client.RootGoitPath, cleanedArg, client.Idx); err != nil {
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
