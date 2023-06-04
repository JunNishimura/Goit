package store

import (
	"os"
	"path/filepath"

	"github.com/JunNishimura/Goit/internal/sha"
)

type branch struct {
	Name string
	hash sha.SHA1
}

type Refs struct {
	Heads []*branch
}

func NewRefs(rootGoitPath string) (*Refs, error) {
	r := newRefs()
	headsPath := filepath.Join(rootGoitPath, "refs", "heads")
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
