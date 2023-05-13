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

func isIndexNeedUpdated(object *object.Object, filePath string) (bool, error) {
	f, err := os.Open(INDEX_PATH)
	if err != nil {
		return false, fmt.Errorf("fail to open %s: %v", INDEX_PATH, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		blobInfo := strings.Split(scanner.Text(), " ")
		if len(blobInfo) != 2 {
			return false, fmt.Errorf("find invalid blob info %v", blobInfo)
		}
		// if blob which has same path and hash is registered, return false.
		if blobInfo[0] == object.Hash.String() && blobInfo[1] == filePath {
			return false, nil
		}
	}
	return true, nil
}

func updateIndex(object *object.Object, filePath string) error {
	f, err := os.OpenFile(INDEX_PATH, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("fail to open %s: %v", INDEX_PATH, err)
	}

	// store lines of file except the line which has the same filePath
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		blobInfo := strings.Split(scanner.Text(), " ")
		if blobInfo[0] != object.Hash.String() && blobInfo[1] == filePath {
			// skip the same file
			continue
		}
		lines = append(lines, strings.Join(blobInfo, " "))
	}
	lines = append(lines, strings.Join([]string{object.Hash.String(), filePath}, " "))
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

			// check if index needs to be updated
			updateFlag, err := indexClient.IsUpdateNeeded()
			if err != nil {
				return fmt.Errorf("fail to see if index needs to be updated: %v", err)
			}
			if !updateFlag {
				continue
			}

			// update index
			path := []byte(arg) //TODO: update path construction
			if err := indexClient.Update(object.Hash, path); err != nil {
				return fmt.Errorf("fail to update index: %v", err)
			}

			// compress file by zlib
			compData, err := object.CompressBlob()
			if err != nil {
				return fmt.Errorf("fail to compress data: %v", err)
			}

			// save file
			object.Write(compData)
		}

		// delete non-tracking files from index
		// if err := deleteFromIndex(); err != nil {
		// 	return fmt.Errorf("fail to delete from index: %v", err)
		// }

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
