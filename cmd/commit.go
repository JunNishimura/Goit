/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/JunNishimura/Goit/internal/log"
	"github.com/JunNishimura/Goit/internal/object"
	"github.com/JunNishimura/Goit/internal/sha"
	"github.com/JunNishimura/Goit/internal/store"
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

func commit(rootGoitPath string, index *store.Index, head *store.Head, conf *store.Config, refs *store.Refs) error {
	// make and write tree object
	treeObject, err := writeTreeObject(rootGoitPath, index.Entries)
	if err != nil {
		return err
	}

	// make and write commit object
	var data []byte
	branchPath := filepath.Join(rootGoitPath, "refs", "heads", head.Reference)
	branchBytes, err := os.ReadFile(branchPath)
	author := object.NewSign(conf.GetUserName(), conf.GetEmail())
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
	if err := commit.Write(rootGoitPath); err != nil {
		return fmt.Errorf("fail to write commit object: %w", err)
	}

	// create/update branch
	var from sha.SHA1
	if refs.IsBranchExist(head.Reference) {
		// update
		if err := refs.UpdateBranchHash(rootGoitPath, head.Reference, commit.Hash); err != nil {
			return fmt.Errorf("fail to update branch %s: %w", head.Reference, err)
		}
		from = head.Commit.Hash
	} else {
		// create
		if err := refs.AddBranch(rootGoitPath, head.Reference, commit.Hash); err != nil {
			return fmt.Errorf("fail to create branch %s: %w", head.Reference, err)
		}
		from = nil
	}
	// log
	record := log.NewRecord(log.CommitRecord, from, commit.Hash, conf.GetUserName(), conf.GetEmail(), time.Now(), message)
	if err := gLogger.WriteHEAD(record); err != nil {
		return fmt.Errorf("log error: %w", err)
	}
	if err := gLogger.WriteBranch(record, head.Reference); err != nil {
		return fmt.Errorf("log error: %w", err)
	}

	// update HEAD
	if err := head.Update(refs, rootGoitPath, head.Reference); err != nil {
		return fmt.Errorf("fail to update HEAD: %w", err)
	}

	return nil
}

func isCommitNecessary(rootGoitPath string, index *store.Index, commitObj *object.Commit) (bool, error) {
	// get tree object
	treeObject, err := object.GetObject(rootGoitPath, commitObj.Tree)
	if err != nil {
		return false, fmt.Errorf("fail to get tree object: %w", err)
	}

	// get tree
	tree, err := object.NewTree(rootGoitPath, treeObject)
	if err != nil {
		return false, fmt.Errorf("fail to get tree: %w", err)
	}

	// compare index with tree
	diffEntries, err := index.DiffWithTree(tree)
	if err != nil {
		return false, fmt.Errorf("fail to compare index with tree: %w", err)
	}

	return len(diffEntries) != 0, nil
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
			if err := commit(client.RootGoitPath, client.Idx, client.Head, client.Conf, client.Refs); err != nil {
				return err
			}
		} else {
			// compare last commit with index
			isDiff, err := isCommitNecessary(client.RootGoitPath, client.Idx, client.Head.Commit)
			if err != nil {
				return fmt.Errorf("fail to compare last commit with index: %w", err)
			}
			if !isDiff {
				return ErrNothingToCommit
			}

			// commit
			if err := commit(client.RootGoitPath, client.Idx, client.Head, client.Conf, client.Refs); err != nil {
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
