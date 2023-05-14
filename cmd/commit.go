/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/JunNishimura/Goit/object"
	"github.com/JunNishimura/Goit/sha"
	"github.com/spf13/cobra"
)

var (
	message string
)

func commit() error {
	// make and write tree object
	treeObject, err := object.WriteTreeObject(indexClient.Entries)
	if err != nil {
		return err
	}

	// make and write commit object
	var data []byte
	branchPath := filepath.Join(".goit", "refs", "heads", "main")
	branchBytes, err := os.ReadFile(branchPath)
	if err != nil {
		// no branch means that this is the initial commit
		data = []byte(fmt.Sprintf("tree %s\n\n%s\n", treeObject.Hash, message))
	} else {
		parentHash := string(branchBytes)
		data = []byte(fmt.Sprintf("tree %s\nparent %s\n\n%s\n", treeObject.Hash, parentHash, message))
	}
	commitObject := object.NewObject(object.CommitObject, data)
	commit, err := object.NewCommit(commitObject)
	if err != nil {
		return fmt.Errorf("fail to make commit object: %v", err)
	}
	if err := commit.Write(); err != nil {
		return fmt.Errorf("fail to write commit object: %v", err)
	}

	// update branch
	if err := commit.UpdateBranch(); err != nil {
		return fmt.Errorf("fail to make new branch: %v", err)
	}

	return nil
}

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

			// commit
			if err := commit(); err != nil {
				return err
			}
		} else {
			// get last commit object
			branchPath := filepath.Join(".goit", "refs", "heads", "main")
			hashBytes, err := ioutil.ReadFile(branchPath)
			if err != nil {
				return fmt.Errorf("fail to read %s: %v", branchPath, err)
			}
			hashString := string(hashBytes)
			lastCommitHash, err := sha.ReadHash(hashString)
			if err != nil {
				return fmt.Errorf("fail to decode hash string: %v", err)
			}
			lastCommitObject, err := object.GetObject(lastCommitHash)
			if err != nil {
				return fmt.Errorf("fail to get last commit object: %v", err)
			}

			// get last commit
			lastCommit, err := object.NewCommit(lastCommitObject)
			if err != nil {
				return fmt.Errorf("fail to get last commit: %v", err)
			}

			// compare last commit with index
			isCommitNecessary, err := lastCommit.IsCommitNecessary(indexClient)
			if err != nil {
				return fmt.Errorf("fail to compare last commit with index: %v", err)
			}
			if !isCommitNecessary {
				return errors.New("nothing to commit")
			}

			// commit
			if err := commit(); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringVarP(&message, "message", "m", "", "commit message")
}
