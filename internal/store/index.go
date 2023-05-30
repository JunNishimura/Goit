package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/JunNishimura/Goit/internal/sha"
)

const (
	newEntryFlag = -1
)

type Entry struct {
	Hash       sha.SHA1
	NameLength uint16
	Path       []byte
}

func NewEntry(hash sha.SHA1, path []byte) *Entry {
	return &Entry{
		Hash:       hash,
		NameLength: uint16(len(path)),
		Path:       path,
	}
}

type Header struct {
	Signature [4]byte
	Version   uint32
	EntryNum  uint32
}

type Index struct {
	Header
	Entries []*Entry // sorted entries
}

func NewIndex(rootGoitPath string) (*Index, error) {
	index := &Index{
		Header: Header{
			Signature: [4]byte{'D', 'I', 'R', 'C'},
			Version:   uint32(1),
			EntryNum:  uint32(0),
		},
	}
	indexPath := filepath.Join(rootGoitPath, "index")
	if _, err := os.Stat(indexPath); !os.IsNotExist(err) {
		if err := index.read(rootGoitPath); err != nil {
			return nil, fmt.Errorf("fail to read index: %w", err)
		}
	}
	return index, nil
}

// return the position of entry, entry, and flag to tell the entry is found or not
func (idx *Index) GetEntry(path []byte) (int, *Entry, bool) {
	if idx.EntryNum == 0 {
		return newEntryFlag, nil, false
	}

	left := 0
	right := int(idx.EntryNum)
	for {
		middle := (left + right) / 2
		entry := idx.Entries[middle]
		if string(entry.Path) == string(path) {
			return middle, entry, true
		} else if string(entry.Path) < string(path) {
			left = middle + 1
		} else {
			right = middle
		}

		if right-left < 1 {
			break
		}
	}

	return newEntryFlag, nil, false
}

func (idx *Index) Update(rootGoitPath string, hash sha.SHA1, path []byte) (bool, error) {
	pos, gotEntry, isFound := idx.GetEntry(path)
	if isFound && string(gotEntry.Hash) == string(hash) && string(gotEntry.Path) == string(path) {
		return false, nil
	}

	// add new entry and update index entries
	entry := NewEntry(hash, path)
	if pos != newEntryFlag {
		// remove existing entry
		idx.Entries = append(idx.Entries[:pos], idx.Entries[pos+1:]...)
	}
	idx.Entries = append(idx.Entries, entry)
	idx.EntryNum = uint32(len(idx.Entries))
	sort.Slice(idx.Entries, func(i, j int) bool { return string(idx.Entries[i].Path) < string(idx.Entries[j].Path) })

	if err := idx.write(rootGoitPath); err != nil {
		return false, err
	}

	return true, nil
}

func (idx *Index) DeleteEntry(rootGoitPath string, entry *Entry) error {
	pos, _, isFound := idx.GetEntry(entry.Path)
	if !isFound {
		return fmt.Errorf("'%s' is not registered in index, so fail to delete", entry.Path)
	}

	// delete target entry
	idx.Entries = append(idx.Entries[:pos], idx.Entries[pos+1:]...)

	// reset entry num
	idx.EntryNum = uint32(len(idx.Entries))

	// write index
	if err := idx.write(rootGoitPath); err != nil {
		return err
	}

	return nil
}

func (idx *Index) DeleteUntrackedFiles(rootGoitPath string) error {
	var trackedEntries []*Entry
	for _, entry := range idx.Entries {
		if _, err := os.Stat(string(entry.Path)); !os.IsNotExist(err) {
			trackedEntries = append(trackedEntries, entry)
		}
	}

	// no need to delete
	if len(trackedEntries) == int(idx.EntryNum) {
		return nil
	}

	// need to update index
	idx.Entries = trackedEntries
	idx.EntryNum = uint32(len(idx.Entries))
	if err := idx.write(rootGoitPath); err != nil {
		return err
	}

	return nil
}

func (idx *Index) read(rootGoitPath string) error {
	// read index
	indexPath := filepath.Join(rootGoitPath, "index")
	b, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("fail to read index: %w", err)
	}

	// make bytes reader
	buf := bytes.NewReader(b)

	// fixed length decoding
	err = binary.Read(buf, binary.BigEndian, &idx.Header)
	if err != nil {
		return fmt.Errorf("fail to read index header: %w", err)
	}

	// variable length decoding
	for i := 0; i < int(idx.EntryNum); i++ {
		// read hash
		hash := make(sha.SHA1, 20)
		err = binary.Read(buf, binary.BigEndian, &hash)
		if err != nil {
			return fmt.Errorf("fail to read hash from index: %w", err)
		}

		// read file name length
		var nameLength uint16
		err = binary.Read(buf, binary.BigEndian, &nameLength)
		if err != nil {
			return fmt.Errorf("fail to read file name length from index: %w", err)
		}

		// read file path
		path := make([]byte, nameLength)
		err = binary.Read(buf, binary.BigEndian, &path)
		if err != nil {
			return fmt.Errorf("fail to read path from index: %w", err)
		}

		entry := NewEntry(hash, path)
		idx.Entries = append(idx.Entries, entry)
	}

	return nil
}

func (idx *Index) write(rootGoitPath string) error {
	indexPath := filepath.Join(rootGoitPath, "index")
	f, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("fail to create .goit/index: %w", err)
	}
	defer f.Close()

	// fixed length encoding
	if err := binary.Write(f, binary.BigEndian, &idx.Header); err != nil {
		return fmt.Errorf("fail to write fixed-length encoding: %w", err)
	}

	// variable length encoding
	var data []byte
	for _, entry := range idx.Entries {
		bNameLength := make([]byte, 2)
		binary.BigEndian.PutUint16(bNameLength, entry.NameLength)
		data = append(data, entry.Hash...)
		data = append(data, bNameLength...)
		data = append(data, entry.Path...)
	}
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("fail to write variable-length encoding: %w", err)
	}

	return nil
}
