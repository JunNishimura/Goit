package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/JunNishimura/Goit/internal/object"
	"github.com/JunNishimura/Goit/internal/sha"
)

type Head struct {
	Reference string
	Commit    *object.Commit
}

var (
	headRegexp     = regexp.MustCompile("ref: refs/heads/.+")
	ErrInvalidHead = errors.New("error: invalid HEAD format")
	ErrIOHandling  = errors.New("IO handling error")
)

func getHeadCommit(branch, rootGoitPath string) (*object.Commit, error) {
	branchPath := filepath.Join(rootGoitPath, "refs", "heads", branch)
	hashBytes, err := os.ReadFile(branchPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrIOHandling, branchPath)
	}
	hashString := string(hashBytes)
	hash, err := sha.ReadHash(hashString)
	if err != nil {
		return nil, fmt.Errorf("fail to decode hash string: %w", err)
	}
	commitObject, err := object.GetObject(rootGoitPath, hash)
	if err != nil {
		return nil, fmt.Errorf("fail to get last commit object: %w", err)
	}
	commit, err := object.NewCommit(commitObject)
	if err != nil {
		return nil, fmt.Errorf("fail to get last commit: %w", err)
	}
	return commit, nil
}

func NewHead(rootGoitPath string) (*Head, error) {
	head := newHead()

	headPath := filepath.Join(rootGoitPath, "HEAD")
	if _, err := os.Stat(headPath); !os.IsNotExist(err) {
		// get branch
		headByte, err := os.ReadFile(headPath)
		if err != nil {
			return nil, fmt.Errorf("fail to read file: %s", headPath)
		}
		headString := string(headByte)
		if ok := headRegexp.MatchString(headString); !ok {
			return nil, ErrInvalidHead
		}
		headSplit := strings.Split(headString, ": ")
		slashSplit := strings.Split(headSplit[1], "/")
		branch := slashSplit[len(slashSplit)-1]
		head.Reference = branch

		// get commit from branch
		branchPath := filepath.Join(rootGoitPath, "refs", "heads", branch)
		if _, err := os.Stat(branchPath); os.IsNotExist(err) {
			return head, nil
		}
		commit, err := getHeadCommit(branch, rootGoitPath)
		if err != nil {
			return nil, ErrInvalidHead
		}
		head.Commit = commit

		return head, nil
	}

	return head, nil
}

func newHead() *Head {
	return &Head{}
}
