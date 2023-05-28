/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JunNishimura/Goit/internal/sha"
	"github.com/spf13/cobra"
)

func updateReference(refPath string, hash sha.SHA1) error {
	f, err := os.Create(refPath)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, refPath)
	}
	defer f.Close()

	_, err = f.WriteString(hash.String())
	if err != nil {
		return fmt.Errorf("fail to write hash(%s) to %s", hash.String(), refPath)
	}

	return nil
}

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
		refPath := filepath.Join(client.RootGoitPath, args[0])

		// hash validation
		hashString := args[1]
		if len(hashString) != 40 {
			return ErrInvalidHash
		}
		hashPath := filepath.Join(client.RootGoitPath, "objects", hashString[:2], hashString[2:])
		if _, err := os.Stat(hashPath); os.IsNotExist(err) {
			return fmt.Errorf("fatal: trying to write ref '%s' with nonexistent object %s", refPath, hashString)
		}
		hash, err := sha.ReadHash(hashString)
		if err != nil {
			return ErrInvalidHash
		}

		if err := updateReference(refPath, hash); err != nil {
			return fmt.Errorf("fail to update reference %s: %w", refPath, err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateRefCmd)
}
