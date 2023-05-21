/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func revParse(refNames ...string) error {
	for _, refName := range refNames {
		if strings.ToLower(refName) == "head" {
			refPath := filepath.Join(client.RootGoitPath, "refs", "heads", string(client.Head))
			hashBytes, err := ioutil.ReadFile(refPath)
			if err != nil {
				return fmt.Errorf(`fatal: ambiguous argument '%s': unknown revision or path not in the working tree`, refName)
			}
			hashString := string(hashBytes)
			fmt.Println(hashString)
		} else {
			refPath := filepath.Join(client.RootGoitPath, "refs", "heads", refName)
			hashBytes, err := ioutil.ReadFile(refPath)
			if err != nil {
				return fmt.Errorf(`fatal: ambiguous argument '%s': unknown revision or path not in the working tree`, refName)
			}
			hashString := string(hashBytes)
			fmt.Println(hashString)
		}
	}
	return nil
}

// revParseCmd represents the revParse command
var revParseCmd = &cobra.Command{
	Use:   "rev-parse",
	Short: "pick out and massage parameters",
	Long:  "pick out and massage parameters",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := revParse(args...); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(revParseCmd)
}
