/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JunNishimura/Goit/internal/file"
	"github.com/spf13/cobra"
)

var (
	rFlag bool
)

func removeFromWorkingTree(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("fail to delete %s from the working tree: %w", path, err)
		}
	}

	return nil
}

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

			_, _, isRegistered := client.Idx.GetEntry([]byte(cleanedArg))
			isRegisteredAsDir := client.Idx.IsRegisteredAsDirectory(cleanedArg)

			if !(isRegistered || isRegisteredAsDir) {
				return fmt.Errorf("fatal: pathspec '%s' did not match any files", arg)
			}
		}

		// remove file from working tree and index
		for _, arg := range args {
			// if the arg is directory
			if f, err := os.Stat(arg); !os.IsNotExist(err) && f.IsDir() {
				// get file paths under directory
				absPath, err := filepath.Abs(arg)
				if err != nil {
					return fmt.Errorf("fail to convert %s to abs path: %w", arg, err)
				}
				filePaths, err := file.GetFilePathsUnderDirectory(absPath)
				if err != nil {
					return fmt.Errorf("fail to get file paths under directory: %w", err)
				}

				// filePaths are defined as abs paths
				// so, translate them to rel paths
				var relPaths []string
				curPath, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("fail to get current directory: %w", err)
				}
				for _, filePath := range filePaths {
					relPath, err := filepath.Rel(curPath, filePath)
					if err != nil {
						return fmt.Errorf("fail to get relative path: %w", err)
					}
					cleanedRelPath := strings.ReplaceAll(relPath, `\`, "/")
					relPaths = append(relPaths, cleanedRelPath)
				}

				// remove
				for _, relPath := range relPaths {
					// remove from the working tree
					if err := removeFromWorkingTree(relPath); err != nil {
						return err
					}

					// remove from the index
					if err := client.Idx.DeleteEntry(client.RootGoitPath, []byte(relPath)); err != nil {
						return fmt.Errorf("fail to delete '%s' from the index: %w", relPath, err)
					}
				}
			} else {
				cleanedArg := filepath.Clean(arg)
				cleanedArg = strings.ReplaceAll(cleanedArg, `\`, "/")

				// remove from the working tree
				if err := removeFromWorkingTree(cleanedArg); err != nil {
					return err
				}

				// remove from the index
				if err := client.Idx.DeleteEntry(client.RootGoitPath, []byte(cleanedArg)); err != nil {
					return fmt.Errorf("fail to delete '%s' from the index: %w", cleanedArg, err)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	rmCmd.Flags().BoolVarP(&rFlag, "rec", "r", false, "allow recursive removal")
}
