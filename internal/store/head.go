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

func (h *Head) Update(refs *Refs, rootGoitPath, newRef string) error {
	// check if branch exists
	n := refs.getBranchPos(newRef)
	if n == NewBranchFlag {
		return fmt.Errorf("branch %s does not exist", newRef)
	}

	headPath := filepath.Join(rootGoitPath, "HEAD")
	if _, err := os.Stat(headPath); os.IsNotExist(err) {
		return errors.New("fail to find HEAD, cannot update")
	}
	f, err := os.Create(headPath)
	if err != nil {
		return fmt.Errorf("fail to create HEAD: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("ref: refs/heads/%s", newRef)); err != nil {
		return fmt.Errorf("fail to write HEAD: %w", err)
	}

	h.Reference = newRef

	// get commit from branch
	branchPath := filepath.Join(rootGoitPath, "refs", "heads", newRef)
	if _, err := os.Stat(branchPath); os.IsNotExist(err) {
		return fmt.Errorf("fail to find branch %s: %w", newRef, err)
	}
	commit, err := getHeadCommit(newRef, rootGoitPath)
	if err != nil {
		return ErrInvalidHead
	}
	h.Commit = commit

	return nil
}

// reset Head to the specified state by hash
// This method does not change Head.Reference, just change Commit
func (h *Head) Reset(rootGoitPath string, refs *Refs, hash sha.SHA1) error {
	// write branch hash
	if err := refs.UpdateBranchHash(rootGoitPath, h.Reference, hash); err != nil {
		return fmt.Errorf("fail to update branch hash: %w", err)
	}

	// get commit object
	commitObject, err := object.GetObject(rootGoitPath, hash)
	if err != nil {
		return fmt.Errorf("fail to get commit object: %w", err)
	}

	// get commit
	commit, err := object.NewCommit(commitObject)
	if err != nil {
		return fmt.Errorf("fail to get commit: %w", err)
	}

	// update commit
	h.Commit = commit

	return nil
}
