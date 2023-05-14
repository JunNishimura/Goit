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
			commit, err := object.NewCommit(commitObject)
			if err != nil {
				return fmt.Errorf("fail to make commit object: %v", err)
			}
			if err := commit.Write(); err != nil {
				return fmt.Errorf("fail to write commit object: %v", err)
			}

			// make new branch
			if err := commit.UpdateBranch(); err != nil {
				return fmt.Errorf("fail to make new branch: %v", err)
			}
		} else {
			// branchPath := filepath.Join(".goit", "refs", "heads", "main")
			// branchBytes, err := ioutil.ReadFile(branchPath)
			// if err != nil {
			// 	return fmt.Errorf("fail to read %s: %v", branchPath, err)
			// }

			// // get last commit object
			// lastCommitHash, err := hex.DecodeString(string(branchBytes))
			// if err != nil {
			// 	return fmt.Errorf("fail to decode hash string: %v", err)
			// }
			// lastCommitObject, err := object.GetObject(lastCommitHash)
			// if err != nil {
			// 	return fmt.Errorf("fail to get last commit object: %v", err)
			// }

			// lastCommitObjPath := filepath.Join(".goit", "objects", lastCommitHash[:2], lastCommitHash[2:])
			// f, err := os.Open(lastCommitObjPath)
			// if err != nil {
			// 	return fmt.Errorf("fail to open %s: %v", lastCommitObjPath, err)
			// }
			// defer f.Close()
			// zr, err := zlib.NewReader(f)
			// if err != nil {
			// 	return fmt.Errorf("fail to zlib.NewReader: %v", err)
			// }
			// defer zr.Close()

			// treeBytes, err := ioutil.ReadAll()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringVarP(&message, "message", "m", "", "commit message")
}
