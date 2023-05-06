/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JunNishimura/Goit/hash"
	"github.com/JunNishimura/Goit/object"
	"github.com/JunNishimura/Goit/util"
	"github.com/spf13/cobra"
)

const (
	INDEX_PATH = ".goit/index"
)

func isIndexNeedUpdated(hash, filePath string) (bool, error) {
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
		if blobInfo[0] == hash && blobInfo[1] == filePath {
			return false, nil
		}
	}
	return true, nil
}

func updateIndex(hash, filePath string) error {
	f, err := os.OpenFile(INDEX_PATH, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("fail to open %s: %v", INDEX_PATH, err)
	}

	// store lines of file except the line which has the same filePath
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		blobInfo := strings.Split(scanner.Text(), " ")
		if blobInfo[0] != hash && blobInfo[1] == filePath {
			// skip the same file
			continue
		}
		lines = append(lines, strings.Join(blobInfo, " "))
	}
	lines = append(lines, strings.Join([]string{hash, filePath}, " "))
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

		// make index file if index file is not found
		if _, err := os.Stat(INDEX_PATH); os.IsNotExist(err) {
			_, err := os.Create(INDEX_PATH)
			if err != nil {
				return fmt.Errorf("fail to make %s: %v", INDEX_PATH, err)
			}
		}

		for _, arg := range args {
			// make object source which is input of hash and zlib
			objSource, err := util.CreateObjectSource(arg, object.BLOB_TYPE)
			if err != nil {
				return fmt.Errorf("fail to generate object source: %v", err)
			}

			// make sha1 hash
			hash := hash.StringToHash(objSource)
			if len(hash) != 40 {
				return errors.New("fail to generate hash")
			}

			// check if index needs to be updated
			indexUpdateFlag, err := isIndexNeedUpdated(hash, arg)
			if err != nil {
				return fmt.Errorf("fail to see if index needs to be updated: %v", err)
			}
			if !indexUpdateFlag {
				continue
			}

			// update index
			if err := updateIndex(hash, arg); err != nil {
				return fmt.Errorf("fail to update index: %v", err)
			}

			// compress file by zlib
			b, err := util.Compress(objSource)
			if err != nil {
				return fmt.Errorf("fail to compress data: %v", err)
			}

			// save file
			dirPath := filepath.Join(object.OBJ_DIR, hash[:2])
			filePath := filepath.Join(dirPath, hash[2:])
			if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
				return fmt.Errorf("fail to make %s: %v", dirPath, err)
			}
			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("fail to make %s: %v", filePath, err)
			}
			defer func() error {
				err := f.Close()
				if err != nil {
					return fmt.Errorf("fail to close the file: %v", err)
				}
				return nil
			}()
			if _, err := f.Write(b.Bytes()); err != nil {
				return fmt.Errorf("fail to write to %s: %v", filePath, err)
			}
		}

		// delete non-tracking files from index
		if err := deleteFromIndex(); err != nil {
			return fmt.Errorf("fail to delete from index: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
