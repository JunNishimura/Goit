/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize Goit",
	Long:  `This is a command to initialize Goit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// .goitディレクトリの存在確認
		if _, err := os.Stat(".goit"); !os.IsNotExist(err) {
			// 既にinitされている場合はreturn
			return errors.New("Goit is already initialized")
		}

		// .goit作成
		if err := os.Mkdir(".goit", os.ModePerm); err != nil {
			return fmt.Errorf("fail to make directory: %s", err.Error())
		}

		// .goit/objects作成
		if err := os.Mkdir(".goit/objects", os.ModePerm); err != nil {
			return fmt.Errorf("fail to make directory: %s", err.Error())
		}

		// .goit/refs作成
		if err := os.Mkdir(".goit/refs", os.ModePerm); err != nil {
			return fmt.Errorf("fail to make directory: %s", err.Error())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}