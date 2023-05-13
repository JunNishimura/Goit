/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/JunNishimura/Goit/index"
	"github.com/JunNishimura/Goit/object"
	"github.com/spf13/cobra"
)

func makeTreeObject(entries []*index.Entry) *object.Object {
	var dirName string
	var data []byte
	var entryBuf []*index.Entry
	i := 0
	for {
		if i >= len(entries) {
			// if the last entry is in the directory
			if dirName != "" {
				treeObject := makeTreeObject(entryBuf)
				data = append(data, []byte(dirName)...)
				data = append(data, 0x00)
				data = append(data, treeObject.Hash...)
			}
			break
		}

		entry := entries[i]
		slashSplit := strings.SplitN(string(entry.Path), "/", 2)
		if len(slashSplit) == 1 {
			if dirName != "" {
				// make tree object from entryBuf
				treeObject := makeTreeObject(entryBuf)
				data = append(data, []byte(dirName)...)
				data = append(data, 0x00)
				data = append(data, treeObject.Hash...)
				// clear dirName and entryBuf
				dirName = ""
				entryBuf = make([]*index.Entry, 0)
			} else {
				data = append(data, entry.Path...)
				data = append(data, 0x00)
				data = append(data, entry.Hash...)
				i++
			}
		} else {
			if dirName == "" {
				dirName = slashSplit[0]
				newEntry := index.NewEntry(entry.Hash, []byte(slashSplit[1]))
				entryBuf = append(entryBuf, newEntry)
				i++
			} else if dirName != "" && dirName == slashSplit[0] {
				// same dir with prev entry
				newEntry := index.NewEntry(entry.Hash, []byte(slashSplit[1]))
				entryBuf = append(entryBuf, newEntry)
				i++
			} else if dirName != "" && dirName != slashSplit[0] {
				treeObject := makeTreeObject(entryBuf)
				data = append(data, []byte(dirName)...)
				data = append(data, 0x00)
				data = append(data, treeObject.Hash...)
				// clear dirName and entryBuf
				dirName = ""
				entryBuf = make([]*index.Entry, 0)
			}
		}
	}

	// make tree object
	treeObject := object.NewObject(object.TreeObject, data)

	return treeObject
}

// writeTreeCmd represents the writeTree command
var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "write tree object from index",
	Long:  "this is a command to write tree object from index",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !IsGoitInitialized() {
			return errors.New("fatal: not a goit repository: .goit")
		}

		// make treeObject from index
		rootTreeObject := makeTreeObject(indexClient.Entries)

		// write root tree object
		if err := rootTreeObject.Write(); err != nil {
			return fmt.Errorf("fail to write tree object: %v", err)
		}

		// print out tree object hash
		fmt.Printf("%s\n", rootTreeObject.Hash)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
