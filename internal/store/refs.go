package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JunNishimura/Goit/internal/sha"
)

type branch struct {
	Name string
	hash sha.SHA1
}

func newBranch(rootGoitPath, branchName string) (*branch, error) {
	branchPath := filepath.Join(rootGoitPath, "refs", "heads", branchName)
	hashByte, err := os.ReadFile(branchPath)
	if err != nil {
		return nil, err
	}
	hashString := string(hashByte)
	hash, err := sha.ReadHash(hashString)
	if err != nil {
		return nil, err
	}
	return &branch{
		Name: branchName,
		hash: hash,
	}, nil
}

type Refs struct {
	Heads []*branch
}

func NewRefs(rootGoitPath string) (*Refs, error) {
	r := newRefs()
	headsPath := filepath.Join(rootGoitPath, "refs", "heads")
	if _, err := os.Stat(headsPath); os.IsNotExist(err) {
		return r, nil
	}
	files, err := os.ReadDir(headsPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		b, err := newBranch(rootGoitPath, file.Name())
		if err != nil {
			return nil, err
		}
		r.Heads = append(r.Heads, b)
	}
	return r, nil
}

func newRefs() *Refs {
	return &Refs{
		Heads: make([]*branch, 0),
	}
}

func (r *Refs) GetBranch(name string) (*branch, error) {
	for _, b := range r.Heads {
		if b.Name == name {
			return b, nil
		}
	}
	return nil, fmt.Errorf("fail to find '%s' branch", name)
}

func getBranchPos(branches []*branch, name string) (int, bool) {
	for n, branch := range branches {
		if branch.Name == name {
			return n, true
		}
	}
	return -1, false
}

func (r *Refs) DeleteBranch(rootGoitPath, name string) error {
	// delete branch from Refs
	p, isBranchFound := getBranchPos(r.Heads, name)
	if !isBranchFound {
		return fmt.Errorf("branch '%s' not found", name)
	}
	r.Heads = append(r.Heads[:p], r.Heads[p+1:]...)

	// delete branch file
	branchPath := filepath.Join(rootGoitPath, "refs", "heads", name)
	if err := os.Remove(branchPath); err != nil {
		return err
	}

	return nil
}
