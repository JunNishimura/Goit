package store

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JunNishimura/Goit/internal/log"
	"github.com/JunNishimura/Goit/internal/sha"
	"github.com/fatih/color"
)

type LogRecord struct {
	Hash       sha.SHA1
	isHead     bool
	references []string
	recType    log.RecordType
	message    string
}

type Reflog struct {
	records []*LogRecord
}

func NewReflog(rootGoitPath string, head *Head, refs *Refs) (*Reflog, error) {
	reflog := newReflog()
	if err := reflog.load(rootGoitPath, head, refs); err != nil {
		return nil, err
	}

	return reflog, nil
}

func newReflog() *Reflog {
	return &Reflog{
		records: make([]*LogRecord, 0),
	}
}

func (r *Reflog) load(rootGoitPath string, head *Head, refs *Refs) error {
	headPath := filepath.Join(rootGoitPath, "logs", "HEAD")
	f, err := os.Open(headPath)
	if err != nil {
		return fmt.Errorf("fail to open %s: %w", headPath, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		record := &LogRecord{
			references: make([]string, 0),
		}

		text := scanner.Text()

		// extract hash
		sp1 := strings.SplitN(text, " ", 3)
		if len(sp1) != 3 {
			continue
		}
		if sp1[1] == strings.Repeat("0", 40) {
			record.Hash = nil
		} else {
			hash, err := sha.ReadHash(sp1[1])
			if err != nil {
				return fmt.Errorf("fail to read hash %s: %w", sp1[1], err)
			}
			record.Hash = hash

			// references
			if head.Commit.Hash.Compare(hash) {
				record.isHead = true
			}
			branches := refs.getBranchesByHash(hash)
			for _, branch := range branches {
				record.references = append(record.references, color.GreenString(branch.Name))
			}
		}

		// extract recType
		sp2 := strings.Split(sp1[2], "\t")
		if len(sp2) != 2 {
			continue
		}
		sp3 := strings.Split(sp2[1], ": ")
		if len(sp3) != 2 {
			continue
		}
		recType := log.NewRecordType(sp3[0])
		if recType == log.UndefinedRecord {
			continue
		}
		record.recType = recType
		record.message = sp3[1]

		r.records = append(r.records, record)
	}

	return nil
}

func (r *Reflog) Show() {
	for i := range r.records {
		idx := len(r.records) - i - 1
		record := r.records[idx]

		var referenceString string
		if len(record.references) > 0 {
			referenceString = strings.Join(record.references, ", ")
		}
		if record.isHead {
			referenceString = color.BlueString("HEAD -> ") + referenceString
		}

		if referenceString == "" {
			fmt.Printf("%s HEAD@{%d}: %s: %s\n", color.YellowString(record.Hash.String()[:7]), i, record.recType, record.message)
		} else {
			fmt.Printf("%s (%s) HEAD@{%d}: %s: %s\n", color.YellowString(record.Hash.String()[:7]), referenceString, i, record.recType, record.message)
		}
	}
}
