/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	maxCount int
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "print commit log",
	Long:  "this is a command to print commit log",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !IsGoitInitialized() {
			return errors.New("fatal: not a goit repository: .goit")
		}

		// see if committed before
		dirName := filepath.Join(client.RootGoitPath, "refs", "heads")
		files, err := ioutil.ReadDir(dirName)
		if err != nil {
			return fmt.Errorf("fail to read dir %s: %v", dirName, err)
		}
		if len(files) == 0 {
			return fmt.Errorf("fatal: your current branch 'main' does not have any commits yet")
		}

		// get last commit hash
		branchPath := filepath.Join(dirName, "main")
		lastCommitHashBytes, err := ioutil.ReadFile(branchPath)
		if err != nil {
			return fmt.Errorf("fail to read %s: %v", branchPath, err)
		}
		lastCommitHashString := string(lastCommitHashBytes)
		fmt.Println(lastCommitHashString)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().IntVarP(&maxCount, "max-count", "n", 5, "max count of logs to print")
}
