/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/JunNishimura/Goit/internal/diff"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "show changes between commits, commit and working tree, etc",
	Long:  "show changes between commits, commit and working tree, etc",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		d := diff.NewDiff([]rune(args[0]), []rune(args[1]))

		d.Compose()

		fmt.Println(d.EditDistance)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}