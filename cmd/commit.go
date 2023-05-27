/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/JunNishimura/Goit/internal/object"
	"github.com/JunNishimura/Goit/internal/sha"
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
	ErrNothingToCommit = errors.New("nothing to commit, working tree clean")
)

func commit() error {
	// make and write tree object
	treeObject, err := writeTreeObject(client.RootGoitPath, client.Idx.Entries)
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
	commitObject, err := object.NewObject(object.CommitObject, data)
	if err != nil {
		return fmt.Errorf("fail to get new object: %w", err)
	}
	commit, err := object.NewCommit(commitObject)
	if err != nil {
		return fmt.Errorf("fail to make commit object: %w", err)
	}
	if err := commit.Write(client.RootGoitPath); err != nil {
		return fmt.Errorf("fail to write commit object: %w", err)
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
		return false, fmt.Errorf("fail to get tree object: %w", err)
	}

	// get entries from tree object
	entries, err := treeObject.ExtractEntries(client.RootGoitPath, "")
	if err != nil {
		return false, fmt.Errorf("fail to get filepath from tree object: %w", err)
	}

	// compare entries extraceted from tree object with index
	if len(entries) != int(client.Idx.EntryNum) {
		return true, nil
	}
	for i := 0; i < len(entries); i++ {
		if string(entries[i].Path) != string(client.Idx.Entries[i].Path) {
			return true, nil
		}
		if entries[i].Hash.String() != client.Idx.Entries[i].Hash.String() {
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !client.Conf.IsUserSet() {
			return ErrUserNotSetOnConfig
		}

		// see if committed before
		dirName := filepath.Join(client.RootGoitPath, "refs", "heads")
		files, err := os.ReadDir(dirName)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrIOHandling, dirName)
		}

		if len(files) == 0 { // no commit before
			if client.Idx.EntryNum == 0 {
				return ErrNothingToCommit
			}

			// commit
			if err := commit(); err != nil {
				return err
			}
		} else {
			// check if files are deleted
			indexPath := filepath.Join(client.RootGoitPath, "index")
			if err := client.Idx.DeleteUntrackedFiles(indexPath); err != nil {
				return fmt.Errorf("fail to delete untracked files: %w", err)
			}

			// get last commit object
			branchPath := filepath.Join(client.RootGoitPath, "refs", "heads", "main")
			hashBytes, err := os.ReadFile(branchPath)
			if err != nil {
				return fmt.Errorf("%w: %s", ErrIOHandling, branchPath)
			}
			hashString := string(hashBytes)
			lastCommitHash, err := sha.ReadHash(hashString)
			if err != nil {
				return fmt.Errorf("fail to decode hash string: %w", err)
			}
			lastCommitObject, err := object.GetObject(client.RootGoitPath, lastCommitHash)
			if err != nil {
				return fmt.Errorf("fail to get last commit object: %w", err)
			}

			// get last commit
			lastCommit, err := object.NewCommit(lastCommitObject)
			if err != nil {
				return fmt.Errorf("fail to get last commit: %w", err)
			}

			// compare last commit with index
			isCommitNecessary, err := isCommitNecessary(lastCommit)
			if err != nil {
				return fmt.Errorf("fail to compare last commit with index: %w", err)
			}
			if !isCommitNecessary {
				return ErrNothingToCommit
			}

			// commit
			if err := commit(); err != nil {
				return fmt.Errorf("fail to commit: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringVarP(&message, "message", "m", "", "commit message")
}
