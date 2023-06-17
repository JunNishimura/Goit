/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/JunNishimura/Goit/internal/store"
	"github.com/spf13/cobra"
)

var (
	isSoft      bool
	isMixed     bool
	isHard      bool
	resetRegexp = regexp.MustCompile(`HEAD@\{\d\}`)
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset current HEAD to the specified state",
	Long:  "reset current HEAD to the specified state",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if client.RootGoitPath == "" {
			return ErrGoitNotInitialized
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// flag validation
		if isSoft || isHard {
			isMixed = false
		}
		if !((isSoft && !isMixed && !isHard) ||
			(!isSoft && isMixed && !isHard) ||
			(!isSoft && !isMixed && isHard)) {
			return errors.New("invalid flags")
		}

		// args validation
		if !(len(args) == 1 && resetRegexp.MatchString(args[0])) {
			return errors.New("only one argument is acceptible. argument format is 'HEAD@{number}'")
		}

		if isSoft {
			// get log record
			reflog, err := store.NewReflog(client.RootGoitPath, client.Head, client.Refs)
			if err != nil {
				return fmt.Errorf("fail to initialize reflog: %w", err)
			}
			sp := strings.Split(args[0], "HEAD@")[1]
			headNum, err := strconv.Atoi(sp[1 : len(sp)-1])
			if err != nil {
				return fmt.Errorf("fail to convert number '%s': %w", args[0], err)
			}
			logRecord, err := reflog.GetRecord(headNum)
			if err != nil {
				return fmt.Errorf("fail to get log record: %w", err)
			}

			// reset Head
			if err := client.Head.Reset(client.RootGoitPath, client.Refs, logRecord.Hash); err != nil {
				return fmt.Errorf("fail to reset HEAD: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	resetCmd.Flags().BoolVar(&isSoft, "soft", false, "reset HEAD")
	resetCmd.Flags().BoolVar(&isMixed, "mixed", true, "reset HEAD and index")
	resetCmd.Flags().BoolVar(&isHard, "hard", false, "reset HEAD, index and working tree")
}
