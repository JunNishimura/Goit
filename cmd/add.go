/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/JunNishimura/Goit/hash"
	"github.com/JunNishimura/Goit/object"
	"github.com/JunNishimura/Goit/util"
	"github.com/spf13/cobra"
)

func addObject(cmd *cobra.Command, args []string) error {
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
		// make object source which is input of hash and zlib
		objSource, err := util.CreateObjectSource(arg, object.Blob)
		if err != nil {
			return fmt.Errorf("fail to generate object source: %v", err)
		}

		// make sha1 hash
		hash := hash.StringToHash(objSource)
		if len(hash) != 40 {
			return errors.New("fail to generate hash")
		}

		// compress file by zlib
		b, err := util.Compress(objSource)
		if err != nil {
			return fmt.Errorf("fail to compress data: %v", err)
		}

		// save file
		dirPath := strings.Join([]string{object.OBJ_DIR, hash[:2]}, "/")
		filePath := strings.Join([]string{dirPath, hash[2:]}, "/")
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

	return nil
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "register changes to index",
	Long:  `This is a command to register changes to index.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := addObject(cmd, args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
