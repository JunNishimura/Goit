package store

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	directoryRegexp = regexp.MustCompile(`.*\/`)
)

type Ignore struct {
	paths []string
}

func NewIgnore(rootGoitPath string) (*Ignore, error) {
	i := newIgnore()
	if err := i.load(rootGoitPath); err != nil {
		return nil, err
	}
	return i, nil
}

func newIgnore() *Ignore {
	return &Ignore{
		paths: []string{`\.goit/.*`},
	}
}

func (i *Ignore) load(rootGoitPath string) error {
	goitignorePath := filepath.Join(filepath.Dir(rootGoitPath), ".goitignore")
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
		var replacedText string
		if directoryRegexp.MatchString(text) {
			replacedText = fmt.Sprintf("%s.*", text)
		} else {
			replacedText = strings.ReplaceAll(text, ".", `\.`)
			replacedText = strings.ReplaceAll(replacedText, "*", ".*")
		}
		i.paths = append(i.paths, replacedText)
	}

	return nil
}

// return true if the parameter is included in ignore list
func (i *Ignore) IsIncluded(path string) bool {
	target := path
	info, _ := os.Stat(path)
	if info.IsDir() && !directoryRegexp.MatchString(path) {
		target = fmt.Sprintf("%s/", path)
	}
	for _, exFile := range i.paths {
		exRegexp := regexp.MustCompile(exFile)
		if exRegexp.MatchString(target) {
			return true
		}
	}
	return false
}
