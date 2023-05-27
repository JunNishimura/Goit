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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// args validation check
		if len(args) != 2 {
			return ErrInvalidArgs
		}
		dotSplit := strings.Split(args[0], ".")
		if len(dotSplit) != 2 {
			return ErrInvalidArgs
		}

		// get global flag
		isGlobal, err := cmd.Flags().GetBool("global")
		if err != nil {
			return fmt.Errorf("fail to get global falg: %w", err)
		}

		// check the existence of global config file
		var globalConfigPath string
		if isGlobal {
			userHomePath, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("fail to get user home dir: %w", err)
			}
			globalConfigPath = filepath.Join(userHomePath, ".goitconfig")
			if _, err := os.Stat(globalConfigPath); os.IsNotExist(err) {
				if _, err := os.Create(globalConfigPath); err != nil {
					return fmt.Errorf("%w: %s", ErrIOHandling, globalConfigPath)
				}
			}
		}
		// check the existence of local config file
		localConfigPath := filepath.Join(client.RootGoitPath, "config")
		if _, err := os.Stat(localConfigPath); os.IsNotExist(err) {
			// if there is no config file, make it
			if _, err := os.Create(localConfigPath); err != nil {
				return fmt.Errorf("%w: %s", ErrIOHandling, localConfigPath)
			}
		}

		// add to config
		identifier := dotSplit[0]
		key := dotSplit[1]
		value := args[1]
		client.Conf.Add(identifier, key, value, isGlobal)

		// write to config
		if err := client.Conf.Write(localConfigPath, globalConfigPath); err != nil {
			return fmt.Errorf("fail to write config: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().Bool("global", false, "add global setting")
}
