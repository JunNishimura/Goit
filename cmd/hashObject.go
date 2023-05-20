/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/JunNishimura/Goit/object"
	"github.com/spf13/cobra"
)

// hashObjectCmd represents the hashObject command
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "calculate the hash of the file",
	Long:  "calculate the hash of the file",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !isGoitInitialized() {
			return errors.New("fatal: not a goit repository: .goit")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			// check if arg is valid
			f, err := os.Stat(arg)
			if os.IsNotExist(err) {
				return fmt.Errorf(`fatal: Cannot open '%s': No such file`, arg)
			}
			if f.IsDir() {
				return fmt.Errorf(`fatal: '%s' is invalid to make blob object`, arg)
			}

			// get data from file
			data, err := ioutil.ReadFile(arg)
			if err != nil {
				return fmt.Errorf("fail to read file: %v", err)
			}

			// make blob object
			object := object.NewObject(object.BlobObject, data)

			fmt.Printf("%s\n", object.Hash)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(hashObjectCmd)
}
