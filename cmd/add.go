/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/JunNishimura/Goit/object"
	"github.com/spf13/cobra"
)

const (
	INDEX_PATH = ".goit/index"
)

func deleteFromIndex() error {
	f, err := os.Open(INDEX_PATH)
	if err != nil {
		return fmt.Errorf("fail to open %s: %v", INDEX_PATH, err)
	}

	// remove files to be deleted from index
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		blobInfo := strings.Split(scanner.Text(), " ")
		if _, err := os.Stat(blobInfo[1]); !os.IsNotExist(err) {
			lines = append(lines, strings.Join(blobInfo, " "))
		}
	}
	f.Close()

	// rewrite index
	f, err = os.OpenFile(INDEX_PATH, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("fail to open %s: %v", INDEX_PATH, err)
	}
	for _, line := range lines {
		fmt.Fprintln(f, line)
	}
	defer f.Close()

	return nil
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "register changes to index",
	Long:  `This is a command to register changes to index.`,
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
			// make blob object
			object, err := object.NewBlobObject(arg)
			if err != nil {
				return err
			}

			// update index
			path := []byte(arg) //TODO: update path construction
			isUpdated, err := indexClient.Update(object.Hash, path)
			if err != nil {
				return fmt.Errorf("fail to update index: %v", err)
			}
			if !isUpdated {
				continue
			}

			// compress file by zlib
			compData, err := object.CompressBlob()
			if err != nil {
				return fmt.Errorf("fail to compress data: %v", err)
			}

			// save file
			object.Write(compData)
		}

		// delete untracked files from index
		if err := indexClient.DeleteUntrackedFiles(); err != nil {
			return fmt.Errorf("fail to delete untracked files from index: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
