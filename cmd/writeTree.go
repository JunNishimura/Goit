/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/JunNishimura/Goit/internal/object"
	index "github.com/JunNishimura/Goit/internal/store"
	"github.com/spf13/cobra"
)

func writeTreeObject(rootGoitPath string, entries []*index.Entry) (*object.Object, error) {
	var dirName string
	var data []byte
	var entryBuf []*index.Entry
	i := 0
	for {
		if i >= len(entries) {
			// if the last entry is in the directory
			if dirName != "" {
				treeObject, err := writeTreeObject(rootGoitPath, entryBuf)
				if err != nil {
					return nil, err
				}
				data = append(data, []byte(fmt.Sprintf("040000 %s", dirName))...)
				data = append(data, 0x00)
				data = append(data, treeObject.Hash...)
			}
			break
		}

		entry := entries[i]
		slashSplit := strings.SplitN(string(entry.Path), "/", 2)
		if len(slashSplit) == 1 { // if entry is not in sub-directory
			if dirName != "" { // if previous entry is in sub-directory
				// make tree object from entryBuf
				treeObject, err := writeTreeObject(rootGoitPath, entryBuf)
				if err != nil {
					return nil, err
				}
				data = append(data, []byte(fmt.Sprintf("040000 %s", dirName))...)
				data = append(data, 0x00)
				data = append(data, treeObject.Hash...)
				// clear dirName and entryBuf
				dirName = ""
				entryBuf = make([]*index.Entry, 0)
			}
			data = append(data, []byte(fmt.Sprintf("100644 %s", string(entry.Path)))...)
			data = append(data, 0x00)
			data = append(data, entry.Hash...)
		} else { // if entry is in sub-directory
			if dirName == "" { // previous entry is not in sub-directory
				dirName = slashSplit[0] // root sub-directory name e.x) cmd/pkg/main.go -> cmd
				newEntry := index.NewEntry(entry.Hash, []byte(slashSplit[1]))
				entryBuf = append(entryBuf, newEntry)
			} else if dirName != "" && dirName == slashSplit[0] { // previous entry is in sub-directory, and current entry is in the same sub-directory
				newEntry := index.NewEntry(entry.Hash, []byte(slashSplit[1]))
				entryBuf = append(entryBuf, newEntry)
			} else if dirName != "" && dirName != slashSplit[0] { // previous entry is in sub-directory, and current entry is in the different sub-directory
				// make tree object
				treeObject, err := writeTreeObject(rootGoitPath, entryBuf)
				if err != nil {
					return nil, err
				}
				data = append(data, []byte(fmt.Sprintf("040000 %s", dirName))...)
				data = append(data, 0x00)
				data = append(data, treeObject.Hash...)
				// start making tree object for different sub-directory
				dirName = slashSplit[0]
				newEntry := index.NewEntry(entry.Hash, []byte(slashSplit[1]))
				entryBuf = []*index.Entry{newEntry}
			}
		}

		i++
	}

	// make tree object
	treeObject, err := object.NewObject(object.TreeObject, data)
	if err != nil {
		return nil, err
	}

	// write tree object
	if err := treeObject.Write(rootGoitPath); err != nil {
		return nil, err
	}

	return treeObject, nil
}

// writeTreeCmd represents the writeTree command
var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "write tree object from index",
	Long:  "this is a command to write tree object from index",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// make and write treeObject from index
		rootTreeObject, err := writeTreeObject(client.RootGoitPath, client.Idx.Entries)
		if err != nil {
			return fmt.Errorf("fail to write tree object: %w", err)
		}

		// print out tree object hash
		fmt.Printf("%s\n", rootTreeObject.Hash)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
