/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize Goit",
	Long:  "This is a command to initialize Goit.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// check goit initizlied
		if IsGoitInitialized() {
			return errors.New("goit is already initialized")
		}

		// make .goit directory
		if err := os.Mkdir(".goit", os.ModePerm); err != nil {
			return fmt.Errorf("fail to make .goit directory: %v", err)
		}

		// make .goit/objects directory
		objectsDir := filepath.Join(".goit", "objects")
		if err := os.Mkdir(objectsDir, os.ModePerm); err != nil {
			return fmt.Errorf("fail to make %s directory: %v", objectsDir, err)
		}

		// make .goit/refs directory
		refsDir := filepath.Join(".goit", "refs")
		headsDir := filepath.Join(".goit", "refs", "heads")
		tagsDir := filepath.Join(".goit", "refs", "tags")
		if err := os.Mkdir(refsDir, os.ModePerm); err != nil {
			return fmt.Errorf("fail to make %s directory: %v", refsDir, err)
		}
		if err := os.Mkdir(headsDir, os.ModePerm); err != nil {
			return fmt.Errorf("fail to make %s directory: %v", headsDir, err)
		}
		if err := os.Mkdir(tagsDir, os.ModePerm); err != nil {
			return fmt.Errorf("fail to make %s directory: %v", tagsDir, err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
