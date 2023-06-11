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

func commit() error {
	// make and write tree object
	treeObject, err := writeTreeObject(client.RootGoitPath, client.Idx.Entries)
	if err != nil {
		return err
	}

	// make and write commit object
	var data []byte
	branchPath := filepath.Join(client.RootGoitPath, "refs", "heads", client.Head.Reference)
	branchBytes, err := os.ReadFile(branchPath)
	author := object.NewSign(client.Conf.GetUserName(), client.Conf.GetEmail())
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

	// create/update branch
	if client.Refs.IsBranchExist(client.Head.Reference) {
		// update
		if err := client.Refs.UpdateBranchHash(client.RootGoitPath, client.Head.Reference, commit.Hash); err != nil {
			return fmt.Errorf("fail to update branch %s: %w", client.Head.Reference, err)
		}
	} else {
		// create
		if err := client.Refs.AddBranch(client.RootGoitPath, client.Head.Reference, commit.Hash); err != nil {
			return fmt.Errorf("fail to create branch %s: %w", client.Head.Reference, err)
		}
		if err := client.Head.Update(client.Refs, client.RootGoitPath, client.Head.Reference); err != nil {
			return fmt.Errorf("fail to update HEAD: %w", err)
		}
	}

	return nil
}

func getEntriesFromTree(rootName string, nodes []*object.Node) ([]*store.Entry, error) {
	var entries []*store.Entry

	for _, node := range nodes {
		if len(node.Children) == 0 {
			var entryName string
			if rootName == "" {
				entryName = node.Name
			} else {
				entryName = fmt.Sprintf("%s/%s", rootName, node.Name)
			}
			newEntry := &store.Entry{
				Hash:       node.Hash,
				NameLength: uint16(len(entryName)),
				Path:       []byte(entryName),
			}
			entries = append(entries, newEntry)
		} else {
			var newRootName string
			if rootName == "" {
				newRootName = node.Name
			} else {
				newRootName = fmt.Sprintf("%s/%s", rootName, node.Name)
			}
			childEntries, err := getEntriesFromTree(newRootName, node.Children)
			if err != nil {
				return nil, err
			}
			entries = append(entries, childEntries...)
		}
	}

	return entries, nil
}

func isIndexDifferentFromTree(index *store.Index, tree *object.Tree) (bool, error) {
	rootName := ""
	gotEntries, err := getEntriesFromTree(rootName, tree.Children)
	if err != nil {
		return false, err
	}

	if len(gotEntries) != int(index.EntryNum) {
		return true, nil
	}
	for i := 0; i < len(gotEntries); i++ {
		if string(gotEntries[i].Path) != string(index.Entries[i].Path) {
			return true, nil
		}
		if gotEntries[i].Hash.String() != index.Entries[i].Hash.String() {
			return true, nil
		}
	}
	return false, nil
}

func isCommitNecessary(commitObj *object.Commit) (bool, error) {
	// get tree object
	treeObject, err := object.GetObject(client.RootGoitPath, commitObj.Tree)
	if err != nil {
		return false, fmt.Errorf("fail to get tree object: %w", err)
	}

	// get tree
	tree, err := object.NewTree(client.RootGoitPath, treeObject)
	if err != nil {
		return false, fmt.Errorf("fail to get tree: %w", err)
	}

	// compare index with tree
	isDiff, err := isIndexDifferentFromTree(client.Idx, tree)
	if err != nil {
		return false, fmt.Errorf("fail to compare index with tree: %w", err)
	}

	return isDiff, nil
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
			// compare last commit with index
			isCommitNecessary, err := isCommitNecessary(client.Head.Commit)
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
