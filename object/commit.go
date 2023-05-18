package object

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JunNishimura/Goit/sha"
	"github.com/JunNishimura/Goit/store"
)

type Commit struct {
	*Object
	Tree    sha.SHA1
	Parents []sha.SHA1
	Message string
}

func NewCommit(o *Object) (*Commit, error) {
	if o.Type != CommitObject {
		return nil, ErrNotCommitObject
	}

	commit := &Commit{
		Object: o,
	}

	buf := bytes.NewReader(o.Data)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		text := scanner.Text()
		splitText := strings.SplitN(text, " ", 2)
		if len(splitText) != 2 {
			break
		}

		lineType := splitText[0]
		body := splitText[1]

		switch lineType {
		case "tree":
			hash, err := sha.ReadHash(body)
			if err != nil {
				return nil, err
			}
			commit.Tree = hash
		case "parent":
			hash, err := sha.ReadHash(body)
			if err != nil {
				return nil, err
			}
			commit.Parents = append(commit.Parents, hash)
		}
	}

	message := make([]string, 0)
	for scanner.Scan() {
		message = append(message, scanner.Text())
	}
	commit.Message = strings.Join(message, "\n")

	return commit, nil
}

func (c *Commit) UpdateBranch() error {
	filePath := filepath.Join(".goit", "refs", "heads", "main")
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("fail to make %s: %v", filePath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(c.Hash.String()); err != nil {
		return fmt.Errorf("fail to write hash to %s: %v", filePath, err)
	}

	return nil
}

func (c *Commit) IsCommitNecessary(idx *store.Index) (bool, error) {
	treeObject, err := GetObject(c.Tree)
	if err != nil {
		return false, fmt.Errorf("fail to get tree object: %v", err)
	}

	// get entries from tree object
	rootDir := ""
	paths, err := treeObject.extractFilePaths(rootDir)
	if err != nil {
		return false, fmt.Errorf("fail to get entries from tree object: %v", err)
	}

	// compare entries extraceted from tree object with index
	if len(paths) != int(idx.EntryNum) {
		return true, nil
	}
	for i := 0; i < len(paths); i++ {
		if paths[i] != string(idx.Entries[i].Path) {
			return true, nil
		}
	}
	return false, nil
}
