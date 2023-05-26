/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return ErrInvalidArgs
		}
		dotSplit := strings.Split(args[0], ".")
		if len(dotSplit) != 2 {
			return ErrInvalidArgs
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// check the existence of config file
		configPath := filepath.Join(client.RootGoitPath, "config")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// if there is no config file, make it
			if _, err := os.Create(configPath); err != nil {
				return fmt.Errorf("%w: %s", ErrIOHandling, configPath)
			}
		}

		// add to config
		dotSplit := strings.Split(args[0], ".")
		identifier := dotSplit[0]
		key := dotSplit[1]
		value := args[1]
		client.Conf.Add(identifier, key, value)

		// write to config
		if err := client.Conf.Write(client.RootGoitPath); err != nil {
			return fmt.Errorf("fail to write config: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
