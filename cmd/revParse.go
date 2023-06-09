/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JunNishimura/Goit/internal/store"
	"github.com/spf13/cobra"
)

func revParse(rootGoitPath string, head *store.Head, refNames ...string) error {
	for _, refName := range refNames {
		var refPath string
		if strings.ToLower(refName) == "head" {
			refPath = filepath.Join(rootGoitPath, "refs", "heads", head.Reference)
		} else {
			refPath = filepath.Join(rootGoitPath, "refs", "heads", refName)
		}
		hashBytes, err := os.ReadFile(refPath)
		if err != nil {
			return fmt.Errorf(`fatal: ambiguous argument '%s': unknown revision or path not in the working tree`, refName)
		}
		hashString := string(hashBytes)
		fmt.Println(hashString)
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
		if err := revParse(client.RootGoitPath, client.Head, args...); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(revParseCmd)
}
