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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath != "" {
			return errors.New("goit is already initialized")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// make .goit directory
		curPath, err := os.Getwd()
		if err != nil {
			return errors.New("fail to get current path")
		}
		goitDir := filepath.Join(curPath, ".goit")
		if err := os.Mkdir(goitDir, os.ModePerm); err != nil {
			return fmt.Errorf("%w: %s", ErrIOHandling, goitDir)
		}

		// make .goit/config file
		configFile := filepath.Join(goitDir, "config")
		if _, err := os.Create(configFile); err != nil {
			return fmt.Errorf("%w: %s", ErrIOHandling, configFile)
		}

		// make .goit/HEAD file and write main branch
		headFile := filepath.Join(goitDir, "HEAD")
		f, err := os.Create(headFile)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrIOHandling, headFile)
		}
		defer f.Close()
		// set 'main' as default branch
		if _, err := f.WriteString("ref: refs/heads/main"); err != nil {
			return fmt.Errorf("%w: %s", ErrIOHandling, headFile)
		}

		// make .goit/objects directory
		objectsDir := filepath.Join(goitDir, "objects")
		if err := os.Mkdir(objectsDir, os.ModePerm); err != nil {
			return fmt.Errorf("%w: %s", ErrIOHandling, objectsDir)
		}

		// make .goit/refs directory
		refsDir := filepath.Join(goitDir, "refs")
		if err := os.Mkdir(refsDir, os.ModePerm); err != nil {
			return fmt.Errorf("%w: %s", ErrIOHandling, refsDir)
		}

		// make .goit/refs/heads directory
		headsDir := filepath.Join(refsDir, "heads")
		if err := os.Mkdir(headsDir, os.ModePerm); err != nil {
			return fmt.Errorf("%w: %s", ErrIOHandling, headsDir)
		}

		// make .goit/refs/tags directory
		tagsDir := filepath.Join(refsDir, "tags")
		if err := os.Mkdir(tagsDir, os.ModePerm); err != nil {
			return fmt.Errorf("%w: %s", ErrIOHandling, tagsDir)
		}

		// print out message for initialization success
		fmt.Printf("Initialized empty Goit repository in %s\n", goitDir)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
