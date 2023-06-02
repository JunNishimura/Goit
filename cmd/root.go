/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/JunNishimura/Goit/internal/file"
	"github.com/JunNishimura/Goit/internal/store"
	"github.com/spf13/cobra"
)

var (
	client      *store.Client
	goitVersion = ""
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goit",
	Short: "Git made by Golang",
	Long:  "This is a Git-like CLI tool made by Golang",
	RunE: func(cmd *cobra.Command, args []string) error {
		versionFlag, err := cmd.Flags().GetBool("version")
		if err != nil {
			return fmt.Errorf("fail to get version flag: %w", err)
		}

		if versionFlag {
			fmt.Println(goitVersion)
		}

		return nil
	},
}

func Execute(version string) {
	// set version
	if version == "" {
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			goitVersion = buildInfo.Main.Version
		}
	} else {
		goitVersion = version
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootGoitPath, _ := file.FindGoitRoot(".") // ignore the error since the error is not important
	config, err := store.NewConfig(rootGoitPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	index, err := store.NewIndex(rootGoitPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	head, err := store.NewHead(rootGoitPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client = store.NewClient(config, index, head, rootGoitPath)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "v", false, "Show Goit version")
}
