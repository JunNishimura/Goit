/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/JunNishimura/Goit/object"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize Goit",
	Long:  `This is a command to initialize Goit.`,
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
		if err := os.Mkdir(object.OBJ_DIR, os.ModePerm); err != nil {
			return fmt.Errorf("fail to make .goit/objects directory: %v", err)
		}

		// make .goit/refs directory
		REFS_DIR := ".goit/refs"
		REFS_HEADS_DIR := strings.Join([]string{REFS_DIR, "heads"}, "/")
		REFS_TAGS_DIR := strings.Join([]string{REFS_DIR, "tags"}, "/")
		if err := os.Mkdir(REFS_DIR, os.ModePerm); err != nil {
			return fmt.Errorf("fail to make .goit/refs directory: %v", err)
		}
		if err := os.Mkdir(REFS_HEADS_DIR, os.ModePerm); err != nil {
			return fmt.Errorf("fail to make .goit/refs/heads directory: %v", err)
		}
		if err := os.Mkdir(REFS_TAGS_DIR, os.ModePerm); err != nil {
			return fmt.Errorf("fail to make .goit/refs/tags directory: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
