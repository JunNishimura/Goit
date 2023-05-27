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

func add(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, path)
	}

	// make blob object
	object, err := object.NewObject(object.BlobObject, data)
	if err != nil {
		return fmt.Errorf("fail to get new object: %w", err)
	}

	// get relative path
	curPath, err := os.Getwd()
	if err != nil {
		return err
	}
	relPath, err := filepath.Rel(curPath, path)
	if err != nil {
		return err
	}
	cleanedRelPath := strings.ReplaceAll(relPath, `\`, "/") // replace backslash with slash
	byteRelPath := []byte(cleanedRelPath)

	// update index
	indexPath := filepath.Join(client.RootGoitPath, "index")
	isUpdated, err := client.Idx.Update(indexPath, object.Hash, byteRelPath)
	if err != nil {
		return fmt.Errorf("fail to update index: %w", err)
	}
	if !isUpdated {
		return nil
	}

	// write object to file
	if err := object.Write(client.RootGoitPath); err != nil {
		return fmt.Errorf("fail to write object: %w", err)
	}

	return nil
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "register changes to index",
	Long:  "This is a command to register changes to index.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// args validation check
		if len(args) == 0 {
			return errors.New("nothing specified, nothing added")
		}
		for _, arg := range args {
			if _, err := os.Stat(arg); os.IsNotExist(err) {
				return fmt.Errorf(`path "%s" did not match any files`, arg)
			}
		}

		for _, arg := range args {
			path, err := filepath.Abs(arg)
			if err != nil {
				return fmt.Errorf("fail to convert abs path: %s", arg)
			}

			// directory
			if f, err := os.Stat(arg); !os.IsNotExist(err) && f.IsDir() {
				filePaths, err := file.GetFilePathsUnderDirectory(path)
				if err != nil {
					return fmt.Errorf("fail to get file path under directory: %w", err)
				}
				for _, filePath := range filePaths {
					if err := add(filePath); err != nil {
						return err
					}
				}
			} else {
				if err := add(path); err != nil {
					return err
				}
			}
		}

		// delete untracked files from index
		indexPath := filepath.Join(client.RootGoitPath, "index")
		if err := client.Idx.DeleteUntrackedFiles(indexPath); err != nil {
			return fmt.Errorf("fail to delete untracked files from index: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
