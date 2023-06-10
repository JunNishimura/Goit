/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/JunNishimura/Goit/internal/sha"
	"github.com/spf13/cobra"
)

var (
	branchRegexp = regexp.MustCompile("refs/heads/.+")
)

// updateRefCmd represents the updateRef command
var updateRefCmd = &cobra.Command{
	Use:   "update-ref",
	Short: "update reference",
	Long:  "update reference",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return ErrInvalidArgs
		}

		// get reference path
		if ok := branchRegexp.MatchString(args[0]); !ok {
			return fmt.Errorf("invalid branch path %s", args[0])
		}
		branchSplit := strings.Split(args[0], "/")
		branchName := branchSplit[len(branchSplit)-1]

		// hash validation
		hashString := args[1]
		if len(hashString) != 40 {
			return ErrInvalidHash
		}
		hashPath := filepath.Join(client.RootGoitPath, "objects", hashString[:2], hashString[2:])
		if _, err := os.Stat(hashPath); os.IsNotExist(err) {
			return fmt.Errorf("fatal: trying to write ref '%s' with nonexistent object %s", args[0], hashString)
		}
		newHash, err := sha.ReadHash(hashString)
		if err != nil {
			return ErrInvalidHash
		}

		if err := client.Refs.UpdateBranchHash(client.RootGoitPath, branchName, newHash); err != nil {
			return fmt.Errorf("fail to update reference %s: %w", args[0], err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateRefCmd)
}
