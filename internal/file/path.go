package file

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrGoitRootNotFound = errors.New("not a goit repository (or any of the parent directories): .goit")
)

func FindGoitRoot(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", ErrGoitRootNotFound
	}
	goitPath := filepath.Join(absPath, ".goit")
	if f, err := os.Stat(goitPath); !os.IsNotExist(err) && f.IsDir() {
		return goitPath, nil
	}

	parentPath := filepath.Dir(absPath)
	if parentPath == absPath {
		return "", ErrGoitRootNotFound
	}
	return FindGoitRoot(parentPath)
}