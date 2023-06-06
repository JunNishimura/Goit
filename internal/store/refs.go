package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/JunNishimura/Goit/internal/sha"
)

const (
	NewBranchFlag = -1
)

type branch struct {
	Name string
	hash sha.SHA1
}

func newBranch(name string, hash sha.SHA1) *branch {
	return &branch{
		Name: name,
		hash: hash,
	}
}

func (b *branch) loadHash(rootGoitPath string) error {
	branchPath := filepath.Join(rootGoitPath, "refs", "heads", b.Name)
	hashByte, err := os.ReadFile(branchPath)
	if err != nil {
		return err
	}
	hashString := string(hashByte)
	hash, err := sha.ReadHash(hashString)
	if err != nil {
		return err
	}
	b.hash = hash

	return nil
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
		b := newBranch(file.Name(), nil)
		if err := b.loadHash(rootGoitPath); err != nil {
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

func (r *Refs) AddBranch(rootGoitPath, newBranchName string, newBranchHash sha.SHA1) error {
	// check if branch already exists
	n := r.getBranchPos(newBranchName)
	if n != NewBranchFlag {
		return fmt.Errorf("fatal: a branch named '%s' already exists", newBranchName)
	}

	b := newBranch(newBranchName, newBranchHash)
	r.Heads = append(r.Heads, b)

	// write file
	branchPath := filepath.Join(rootGoitPath, "refs", "heads", newBranchName)
	f, err := os.Create(branchPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(newBranchHash.String()); err != nil {
		return err
	}

	// sort heads
	sort.Slice(r.Heads, func(i, j int) bool { return r.Heads[i].Name < r.Heads[j].Name })

	return nil
}

// return the index of branch in the Refs Heads.
// if not found, return NewBranchFlag which is -1.
func (r *Refs) getBranchPos(branchName string) int {
	for n, branch := range r.Heads {
		if branch.Name == branchName {
			return n
		}
	}
	return NewBranchFlag
}

func (r *Refs) DeleteBranch(rootGoitPath, name string) error {
	// delete branch from Refs
	p := r.getBranchPos(name)
	if p == NewBranchFlag {
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
