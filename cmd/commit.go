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
	message               string
	ErrUserNotSetOnConfig = errors.New(`
			
*** Please tell me who you are.

Run

 goit config user.email "you@example.com"
 goit config user.name "Your name"

to set your account's default identity.

	`)
)

func commit() error {
	// make and write tree object
	treeObject, err := object.WriteTreeObject(client.RootGoitPath, client.Idx.Entries)
	if err != nil {
		return err
	}

	// make and write commit object
	var data []byte
	branchPath := filepath.Join(client.RootGoitPath, "refs", "heads", "main")
	branchBytes, err := os.ReadFile(branchPath)
	author := object.NewSign(client.Conf.Map["user"]["name"], client.Conf.Map["user"]["email"])
	committer := author
	if err != nil {
		// no branch means that this is the initial commit
		data = []byte(fmt.Sprintf("tree %s\nauthor %s\ncommitter %s\n\n%s\n", treeObject.Hash, author, committer, message))
	} else {
		parentHash := string(branchBytes)
		data = []byte(fmt.Sprintf("tree %s\nparent %s\nauthor %s\ncommitter %s\n\n%s\n", treeObject.Hash, parentHash, author, committer, message))
	}
	commitObject := object.NewObject(object.CommitObject, data)
	commit, err := object.NewCommit(commitObject)
	if err != nil {
		return fmt.Errorf("fail to make commit object: %v", err)
	}
	if err := commit.Write(client.RootGoitPath); err != nil {
		return fmt.Errorf("fail to write commit object: %v", err)
	}

	// update branch
	if err := commit.UpdateBranch(branchPath); err != nil {
		return fmt.Errorf("fail to update %s: %w", branchPath, err)
	}

	return nil
}

func isCommitNecessary(commitObj *object.Commit) (bool, error) {
	treeObject, err := object.GetObject(client.RootGoitPath, commitObj.Tree)
	if err != nil {
		return false, fmt.Errorf("fail to get tree object: %v", err)
	}

	// get entries from tree object
	paths, err := treeObject.ExtractFilePaths(client.RootGoitPath, "")
	if err != nil {
		return false, fmt.Errorf("fail to get entries from tree object: %v", err)
	}

	// compare entries extraceted from tree object with index
	if len(paths) != int(client.Idx.EntryNum) {
		return true, nil
	}
	for i := 0; i < len(paths); i++ {
		if paths[i] != string(client.Idx.Entries[i].Path) {
			return true, nil
		}
	}
	return false, nil
}

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "commit",
	Long:  "this is a command to commit",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !client.Conf.IsUserSet() {
			return ErrUserNotSetOnConfig
		}

		// see if committed before
		dirName := filepath.Join(client.RootGoitPath, "refs", "heads")
		files, err := ioutil.ReadDir(dirName)
		if err != nil {
			return fmt.Errorf("fail to read dir %s: %v", dirName, err)
		}

		if len(files) == 0 { // no commit before
			if client.Idx.EntryNum == 0 {
				return errors.New("nothing to commit, working tree clean")
			}

			// commit
			if err := commit(); err != nil {
				return err
			}
		} else {
			// get last commit object
			branchPath := filepath.Join(client.RootGoitPath, "refs", "heads", "main")
			hashBytes, err := ioutil.ReadFile(branchPath)
			if err != nil {
				return fmt.Errorf("fail to read %s: %v", branchPath, err)
			}
			hashString := string(hashBytes)
			lastCommitHash, err := sha.ReadHash(hashString)
			if err != nil {
				return fmt.Errorf("fail to decode hash string: %v", err)
			}
			lastCommitObject, err := object.GetObject(client.RootGoitPath, lastCommitHash)
			if err != nil {
				return fmt.Errorf("fail to get last commit object: %v", err)
			}

			// get last commit
			lastCommit, err := object.NewCommit(lastCommitObject)
			if err != nil {
				return fmt.Errorf("fail to get last commit: %v", err)
			}

			// compare last commit with index
			isCommitNecessary, err := isCommitNecessary(lastCommit)
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
