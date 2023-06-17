package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/JunNishimura/Goit/internal/sha"
	"github.com/fatih/color"
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

func (b *branch) write(rootGoitPath string) error {
	branchPath := filepath.Join(rootGoitPath, "refs", "heads", b.Name)
	f, err := os.Create(branchPath)
	if err != nil {
		return fmt.Errorf("fail to create %s: %w", branchPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(b.hash.String()); err != nil {
		return fmt.Errorf("fail to write hash(%s): %w", b.hash, err)
	}

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
	sort.Slice(r.Heads, func(i, j int) bool { return r.Heads[i].Name < r.Heads[j].Name })
	return r, nil
}

func newRefs() *Refs {
	return &Refs{
		Heads: make([]*branch, 0),
	}
}

func (r *Refs) ListBranches(headBranchName string) {
	for _, b := range r.Heads {
		if b.Name == headBranchName {
			color.Green("* %s", b.Name)
		} else {
			fmt.Println(b.Name)
		}
	}
}

func (r *Refs) IsBranchExist(branchName string) bool {
	p := r.getBranchPos(branchName)
	return p != NewBranchFlag
}

func (r *Refs) AddBranch(rootGoitPath, newBranchName string, newBranchHash sha.SHA1) error {
	// check if branch already exists
	n := r.getBranchPos(newBranchName)
	if n != NewBranchFlag {
		return fmt.Errorf("a branch named '%s' already exists", newBranchName)
	}

	b := newBranch(newBranchName, newBranchHash)
	r.Heads = append(r.Heads, b)

	// write file
	if err := b.write(rootGoitPath); err != nil {
		return fmt.Errorf("fail to write branch: %w", err)
	}

	// sort heads
	sort.Slice(r.Heads, func(i, j int) bool { return r.Heads[i].Name < r.Heads[j].Name })

	return nil
}

func (r *Refs) RenameBranch(rootGoitPath, curBranchName, newBranchName string) error {
	// check if new branch name is not used for other branches
	n := r.getBranchPos(newBranchName)
	if n != NewBranchFlag {
		return fmt.Errorf("branch named '%s' already exists", newBranchName)
	}

	// get current branch
	curNum := r.getBranchPos(curBranchName)
	if curNum == NewBranchFlag {
		return fmt.Errorf("head branch '%s' does not exist", curBranchName)
	}

	// rename branch
	r.Heads[curNum].Name = newBranchName
	sort.Slice(r.Heads, func(i, j int) bool { return r.Heads[i].Name < r.Heads[j].Name })

	// rename file
	oldPath := filepath.Join(rootGoitPath, "refs", "heads", curBranchName)
	newPath := filepath.Join(rootGoitPath, "refs", "heads", newBranchName)
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("fail to rename file: %w", err)
	}

	return nil
}

// return the index of branch in the Refs Heads.
// if not found, return NewBranchFlag which is -1.
func (r *Refs) getBranchPos(branchName string) int {
	if len(r.Heads) == 0 {
		return NewBranchFlag
	}

	// binary search
	left := 0
	right := len(r.Heads)
	for {
		middle := (left + right) / 2
		b := r.Heads[middle]
		if b.Name == branchName {
			return middle
		}
		if b.Name < branchName {
			left = middle + 1
		}
		if b.Name > branchName {
			right = middle
		}

		if right-left < 1 {
			break
		}
	}

	return NewBranchFlag
}

func (r *Refs) getBranchesByHash(hash sha.SHA1) []*branch {
	var branches []*branch
	for _, branch := range r.Heads {
		if branch.hash.Compare(hash) {
			branches = append(branches, branch)
		}
	}

	return branches
}

func (r *Refs) DeleteBranch(rootGoitPath, headBranchName, deleteBranchName string) error {
	// branch validation
	if deleteBranchName == headBranchName {
		return fmt.Errorf("cannot delete current branch '%s'", headBranchName)
	}
	n := r.getBranchPos(deleteBranchName)
	if n == NewBranchFlag {
		return fmt.Errorf("branch '%s' not found", deleteBranchName)
	}
	deleteBranch := r.Heads[n]

	// delete from refs
	r.Heads = append(r.Heads[:n], r.Heads[n+1:]...)

	// delete branch file
	branchPath := filepath.Join(rootGoitPath, "refs", "heads", deleteBranchName)
	if err := os.Remove(branchPath); err != nil {
		return fmt.Errorf("fail to delete branch file: %w", err)
	}

	// print out message
	fmt.Printf("Deleted branch %s (was %s).\n", deleteBranch.Name, deleteBranch.hash.String()[:7])

	return nil
}

func (r *Refs) UpdateBranchHash(rootGoitPath, branchName string, newHash sha.SHA1) error {
	n := r.getBranchPos(branchName)
	if n == NewBranchFlag {
		return fmt.Errorf("branch '%s' does not exist", branchName)
	}

	branch := r.Heads[n]
	branch.hash = newHash

	// write file
	if err := branch.write(rootGoitPath); err != nil {
		return fmt.Errorf("fail to write branch: %w", err)
	}

	return nil
}
