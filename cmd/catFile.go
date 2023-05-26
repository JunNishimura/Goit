/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/JunNishimura/Goit/internal/object"
	"github.com/JunNishimura/Goit/internal/sha"
	"github.com/spf13/cobra"
)

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "cat goit object",
	Long:  "this is a command to show the goit object",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return ErrNotSpecifiedHash
		}
		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return ErrTooManyArgs
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
		// get flags
		typeFlag, err := cmd.Flags().GetBool("type")
		if err != nil {
			return fmt.Errorf("fail to get type flag: %w", err)
		}
		printFlag, err := cmd.Flags().GetBool("print")
		if err != nil {
			return fmt.Errorf("fail to get print flag: %w", err)
		}

		// flag check
		if typeFlag && printFlag {
			return ErrIncompatibleFlag
		}

		// get object from hash
		hash, err := sha.ReadHash(args[0])
		if err != nil {
			return ErrInvalidHash
		}
		obj, err := object.GetObject(client.RootGoitPath, hash)
		if err != nil {
			return ErrInvalidHash
		}

		// print object type
		if typeFlag {
			fmt.Printf("%s\n", obj.Type)
		}

		// print object content
		if printFlag {
			if obj.Type == object.TreeObject {
				// hash is written in hexadecimal in tree object
				dataString, err := obj.ConvertDataToString()
				if err != nil {
					return fmt.Errorf("fail to convert tree data to string: %w", err)
				}

				// format output
				var lines []string
				var prevFileMode, prevFileName string
				scanner := bufio.NewScanner(strings.NewReader(dataString))
				for scanner.Scan() {
					line := scanner.Text()
					if line[:6] == "100644" || line[:6] == "040000" { // first line
						lineSplit := strings.Split(line, " ")
						prevFileMode = lineSplit[0]
						prevFileName = lineSplit[1]
					} else {
						var objType string
						if prevFileMode == "100644" {
							objType = "blob"
						} else if prevFileMode == "040000" {
							objType = "tree"
						}
						hashString := line[:40]
						lines = append(lines, fmt.Sprintf("%s %s %s    %s", prevFileMode, objType, hashString, prevFileName))
						if len(line) > 40 {
							lineSplit := strings.Split(line[40:], " ")
							prevFileMode = lineSplit[0]
							prevFileName = lineSplit[1]
						}
					}
				}

				formattedOutput := strings.Join(lines, "\n")
				fmt.Println(formattedOutput)

			} else {
				// hash is written as string
				fmt.Printf("%s", obj.Data)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)

	catFileCmd.Flags().BoolP("type", "t", false, "print object type")
	catFileCmd.Flags().BoolP("print", "p", false, "print object content")
}
