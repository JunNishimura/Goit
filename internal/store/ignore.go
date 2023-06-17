package store

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	directoryRegexp = regexp.MustCompile(`.*\/`)
)

type Ignore struct {
	file      []string
	directory []string
}

func NewIgnore(rootGoitPath string) (*Ignore, error) {
	i := newIgnore()
	if err := i.load(rootGoitPath); err != nil {
		return nil, err
	}
	fmt.Println(i.file, i.directory)
	return i, nil
}

func newIgnore() *Ignore {
	return &Ignore{
		file:      make([]string, 0),
		directory: []string{".goit/"},
	}
}

func (i *Ignore) load(rootGoitPath string) error {
	goitignorePath := filepath.Join(filepath.Dir(rootGoitPath), ".goitignore")
	fmt.Println(goitignorePath)
	if _, err := os.Stat(goitignorePath); os.IsNotExist(err) {
		return nil
	}
	f, err := os.Open(goitignorePath)
	if err != nil {
		return fmt.Errorf("fail to open %s: %w", goitignorePath, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if directoryRegexp.MatchString(text) {
			i.directory = append(i.directory, text)
		} else {
			i.file = append(i.file, text)
		}
	}

	return nil
}
