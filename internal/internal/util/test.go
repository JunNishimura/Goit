package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrIOHandling = errors.New("IO handling error")
)

// initialize goit under testdata directory
func GoitInit() error {
	// make .goit directory
	rootPath, err := filepath.Abs("../../")
	if err != nil {
		return errors.New("fail to get current path")
	}
	testdataPath := filepath.Join(rootPath, "testdata")
	goitDir := filepath.Join(testdataPath, ".goit")
	if err := os.Mkdir(goitDir, os.ModePerm); err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, goitDir)
	}

	// make .goit/config file
	configFile := filepath.Join(goitDir, "config")
	if _, err := os.Create(configFile); err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, configFile)
	}

	// make .goit/HEAD file and write main branch
	headFile := filepath.Join(goitDir, "HEAD")
	f, err := os.Create(headFile)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, headFile)
	}
	defer f.Close()
	// set 'main' as default branch
	if _, err := f.WriteString("ref: refs/heads/main"); err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, headFile)
	}

	// make .goit/objects directory
	objectsDir := filepath.Join(goitDir, "objects")
	if err := os.Mkdir(objectsDir, os.ModePerm); err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, objectsDir)
	}

	// make .goit/refs directory
	refsDir := filepath.Join(goitDir, "refs")
	if err := os.Mkdir(refsDir, os.ModePerm); err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, refsDir)
	}

	// make .goit/refs/heads directory
	headsDir := filepath.Join(refsDir, "heads")
	if err := os.Mkdir(headsDir, os.ModePerm); err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, headsDir)
	}

	// make .goit/refs/tags directory
	tagsDir := filepath.Join(refsDir, "tags")
	if err := os.Mkdir(tagsDir, os.ModePerm); err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, tagsDir)
	}

	return nil
}
