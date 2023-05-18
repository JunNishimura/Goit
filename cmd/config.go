/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config setting",
	Long:  "this is a command to set config",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !IsGoitInitialized() {
			return errors.New("fatal: not a goit repository: .goit")
		}

		// check the existence of config file
		configPath := filepath.Join(".goit", "config")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// if there is no config file, make it
			if _, err := os.Create(configPath); err != nil {
				return fmt.Errorf("fail to make %s file: %v", configPath, err)
			}
		}

		// check if the arguments are valid
		if len(args) != 2 {
			return ErrInvalidArgs
		}
		dotSplit := strings.Split(args[0], ".")
		if len(dotSplit) != 2 {
			return ErrInvalidArgs
		}

		// add to config
		identifier := dotSplit[0]
		key := dotSplit[1]
		value := args[1]
		conf.Add(identifier, key, value)

		// write to config
		if err := conf.Write(); err != nil {
			return fmt.Errorf("fail to write config: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
