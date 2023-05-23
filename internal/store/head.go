package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Head string

var (
	headRegexp = regexp.MustCompile("ref: refs/heads/.+")
)

func NewHead(rootGoitPath string) (Head, error) {
	headPath := filepath.Join(rootGoitPath, "HEAD")
	if _, err := os.Stat(headPath); !os.IsNotExist(err) {
		headByte, err := os.ReadFile(headPath)
		if err != nil {
			return "", fmt.Errorf("fail to read file: %s", headPath)
		}
		headString := string(headByte)
		if ok := headRegexp.MatchString(headString); !ok {
			return "", errors.New("invalid HEAD")
		}
		headSplit := strings.Split(headString, ": ")
		slashSplit := strings.Split(headSplit[1], "/")
		branch := slashSplit[len(slashSplit)-1]
		return Head(branch), nil
	}

	return "main", nil
}
