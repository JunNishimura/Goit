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

	"github.com/JunNishimura/Goit/internal/object"
	"github.com/spf13/cobra"
)

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
		// check if args are valid
		if len(args) == 0 {
			return errors.New("nothing specified, nothing added")
		}
		for _, arg := range args {
			if _, err := os.Stat(arg); os.IsNotExist(err) {
				return fmt.Errorf(`path "%s" did not match any files`, arg)
			}
		}

		for _, arg := range args {
			// get data from file
			arg = filepath.Clean(arg)               // remove unnecessary slash
			arg = strings.ReplaceAll(arg, `\`, "/") // replace backslash with slash
			data, err := os.ReadFile(arg)
			if err != nil {
				return fmt.Errorf("%w: %s", ErrIOHandling, arg)
			}

			// make blob object
			object, err := object.NewObject(object.BlobObject, data)
			if err != nil {
				return fmt.Errorf("fail to get new object: %w", err)
			}

			// update index
			path := []byte(arg)
			indexPath := filepath.Join(client.RootGoitPath, "index")
			isUpdated, err := client.Idx.Update(indexPath, object.Hash, path)
			if err != nil {
				return fmt.Errorf("fail to update index: %w", err)
			}
			if !isUpdated {
				continue
			}

			// write object to file
			if err := object.Write(client.RootGoitPath); err != nil {
				return fmt.Errorf("fail to write object: %w", err)
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
