/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/JunNishimura/Goit/object"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "register changes to index",
	Long:  "This is a command to register changes to index.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !IsGoitInitialized() {
			return errors.New("fatal: not a goit repository: .goit")
		}

		// check if args are valid
		if len(args) == 0 {
			return errors.New("nothing specified, nothing added")
		}
		for _, arg := range args {
			_, err := os.Stat(arg)
			if os.IsNotExist(err) {
				return fmt.Errorf(`path "%s" did not match any files`, arg)
			}
		}

		for _, arg := range args {
			// check if arg is valid
			f, err := os.Stat(arg)
			if os.IsNotExist(err) {
				return fmt.Errorf(`fatal: Cannot open '%s': No such file`, arg)
			}
			if f.IsDir() {
				return fmt.Errorf(`fatal: '%s' is invalid to make blob object`, arg)
			}

			// get data from file
			data, err := ioutil.ReadFile(arg)
			if err != nil {
				return fmt.Errorf("fail to read file: %v", err)
			}

			// make blob object
			object := object.NewObject(object.BlobObject, data)

			// update index
			path := []byte(arg) //TODO: update path construction
			isUpdated, err := client.Idx.Update(object.Hash, path)
			if err != nil {
				return fmt.Errorf("fail to update index: %v", err)
			}
			if !isUpdated {
				continue
			}

			// write object to file
			object.Write(client.RootDir)
		}

		// delete untracked files from index
		if err := client.Idx.DeleteUntrackedFiles(); err != nil {
			return fmt.Errorf("fail to delete untracked files from index: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
