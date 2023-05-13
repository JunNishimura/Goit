/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/JunNishimura/Goit/object"
	"github.com/spf13/cobra"
)

var (
	message string
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "commit",
	Long:  "this is a command to commit",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !IsGoitInitialized() {
			return errors.New("fatal: not a goit repository: .goit")
		}

		// see if committed before
		dirName := filepath.Join(".goit", "refs", "heads")
		files, err := ioutil.ReadDir(dirName)
		if err != nil {
			return fmt.Errorf("fail to read dir %s: %v", dirName, err)
		}

		if len(files) == 0 { // no commit before
			if indexClient.EntryNum == 0 {
				return errors.New("nothing to commit, working tree clean")
			}

			// make and write tree object
			treeObject := object.MakeTreeObject(indexClient.Entries)
			if err := treeObject.Write(); err != nil {
				return fmt.Errorf("fail to write tree object: %v", err)
			}

			// make and write commit object
			data := []byte(fmt.Sprintf("tree %s\n\n%s\n", treeObject.Hash, message))
			commitObject := object.NewObject(object.CommitObject, data)
			if err := commitObject.Write(); err != nil {
				return fmt.Errorf("fail to write commit object: %v", err)
			}

			// make new branch
			if err := commitObject.UpdateBranch(); err != nil {
				return fmt.Errorf("fail to make new branch: %v", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringVarP(&message, "message", "m", "", "commit message")
}
