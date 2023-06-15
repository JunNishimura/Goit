/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/JunNishimura/Goit/internal/store"
	"github.com/spf13/cobra"
)

// reflogCmd represents the reflog command
var reflogCmd = &cobra.Command{
	Use:   "reflog",
	Short: "manage reference logs",
	Long:  "manage reference logs",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		reflog, err := store.NewReflog(client.RootGoitPath, client.Head, client.Refs)
		if err != nil {
			return fmt.Errorf("fail to get reflog: %w", err)
		}

		reflog.Show()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(reflogCmd)
}
