/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JunNishimura/Goit/internal/object"
	"github.com/JunNishimura/Goit/internal/sha"
	"github.com/spf13/cobra"
)

var (
	maxCount int
)

type WalkFunc func(commit *object.Commit) error

func walkHistory(rootGoitPath string, hash sha.SHA1, walkFunc WalkFunc) error {
	queue := []sha.SHA1{hash}
	visitMap := map[string]struct{}{}

	loopCounter := 0
	for len(queue) > 0 {
		loopCounter++
		if loopCounter > maxCount {
			break
		}

		currentHash := queue[0]
		if _, ok := visitMap[currentHash.String()]; ok {
			queue = queue[1:]
			continue
		}
		visitMap[currentHash.String()] = struct{}{}

		commitObject, err := object.GetObject(rootGoitPath, currentHash)
		if err != nil {
			return err
		}

		commit, err := object.NewCommit(commitObject)
		if err != nil {
			return err
		}

		if err := walkFunc(commit); err != nil {
			return err
		}

		queue = append(queue[1:], commit.Parents...)
	}

	return nil
}

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "print commit log",
	Long:  "this is a command to print commit log",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// see if committed before
		dirName := filepath.Join(client.RootGoitPath, "refs", "heads")
		files, err := os.ReadDir(dirName)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrIOHandling, dirName)
		}
		if len(files) == 0 {
			return fmt.Errorf("fatal: your current branch 'main' does not have any commits yet")
		}

		// print log
		if err := walkHistory(client.RootGoitPath, client.Head.Commit.Hash, func(commit *object.Commit) error {
			fmt.Println(commit)
			return nil
		}); err != nil {
			return fmt.Errorf("fail to log: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().IntVarP(&maxCount, "max-count", "n", 5, "max count of logs to print")
}
