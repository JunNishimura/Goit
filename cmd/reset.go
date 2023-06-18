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
	"time"

	"github.com/JunNishimura/Goit/internal/log"
	"github.com/JunNishimura/Goit/internal/object"
	"github.com/JunNishimura/Goit/internal/store"
	"github.com/spf13/cobra"
)

var (
	isSoft      bool
	isMixed     bool
	isHard      bool
	resetRegexp = regexp.MustCompile(`HEAD@\{\d\}`)
)

func resetHead(arg, rootGoitPath string, logRecord *store.LogRecord, head *store.Head, refs *store.Refs, conf *store.Config) error {
	// reset Head
	prevHeadHash := head.Commit.Hash
	if err := head.Reset(rootGoitPath, refs, logRecord.Hash); err != nil {
		return fmt.Errorf("fail to reset HEAD: %w", err)
	}

	// log
	newRecord := log.NewRecord(log.ResetRecord, prevHeadHash, logRecord.Hash, conf.GetUserName(), conf.GetEmail(), time.Now(), fmt.Sprintf("moving to %s", arg))
	if err := gLogger.WriteHEAD(newRecord); err != nil {
		return fmt.Errorf("log error: %w", err)
	}
	if err := gLogger.WriteBranch(newRecord, head.Reference); err != nil {
		return fmt.Errorf("log error: %w", err)
	}

	return nil
}

func resetIndex(rootGoitPath string, logRecord *store.LogRecord, index *store.Index) error {
	// reset index
	if err := index.Reset(rootGoitPath, logRecord.Hash); err != nil {
		return fmt.Errorf("fail to reset index: %w", err)
	}

	return nil
}

func resetWorkingTree(rootGoitPath string, index *store.Index) error {
	for _, entry := range index.Entries {
		obj, err := object.GetObject(rootGoitPath, entry.Hash)
		if err != nil {
			return fmt.Errorf("fail to get object: %w", err)
		}
		if err := obj.ReflectToWorkingTree(rootGoitPath, string(entry.Path)); err != nil {
			return fmt.Errorf("fail to reflect %s to working directory: %w", string(entry.Path), err)
		}
	}

	return nil
}

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

		// reset HEAD
		if isSoft || isMixed || isHard {
			if err := resetHead(args[0], client.RootGoitPath, logRecord, client.Head, client.Refs, client.Conf); err != nil {
				return fmt.Errorf("fail to reset HEAD: %w", err)
			}
		}

		// reset index
		if isMixed || isHard {
			if err := resetIndex(client.RootGoitPath, logRecord, client.Idx); err != nil {
				return fmt.Errorf("fail to reset index: %w", err)
			}
		}

		// reset working tree
		if isHard {
			if err := resetWorkingTree(client.RootGoitPath, client.Idx); err != nil {
				return fmt.Errorf("fail to reset working tree: %w", err)
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
