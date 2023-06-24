package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/JunNishimura/Goit/internal/object"
	"github.com/JunNishimura/Goit/internal/sha"
)

const (
	diffDelete diffType = iota
	diffNew
	diffModified
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

type diffType int

func (t diffType) String() string {
	switch t {
	case diffDelete:
		return "deleted:"
	case diffNew:
		return "new file:"
	case diffModified:
		return "modified:"
	default:
		return ""
	}
}

type DiffEntry struct {
	Dt    diffType
	Entry *Entry
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
	index := newIndex()
	indexPath := filepath.Join(rootGoitPath, "index")
	if _, err := os.Stat(indexPath); !os.IsNotExist(err) {
		if err := index.read(rootGoitPath); err != nil {
			return nil, fmt.Errorf("fail to read index: %w", err)
		}
	}
	return index, nil
}

func newIndex() *Index {
	return &Index{
		Header: Header{
			Signature: [4]byte{'D', 'I', 'R', 'C'},
			Version:   uint32(1),
			EntryNum:  uint32(0),
		},
	}
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

func (idx *Index) GetEntriesByDirectory(dirName string) []*Entry {
	var entries []*Entry

	dirRegexp := regexp.MustCompile(fmt.Sprintf(`%s\/.+`, dirName))
	for _, entry := range idx.Entries {
		if dirRegexp.Match(entry.Path) {
			entries = append(entries, entry)
		}
	}

	return entries
}

func (idx *Index) IsRegisteredAsDirectory(dirName string) bool {
	if idx.EntryNum == 0 {
		return false
	}

	dirRegexp := regexp.MustCompile(fmt.Sprintf(`%s\/.+`, dirName))

	left := 0
	right := int(idx.EntryNum)
	for {
		middle := (left + right) / 2
		entry := idx.Entries[middle]
		if dirRegexp.MatchString(string(entry.Path)) {
			return true
		} else if string(entry.Path) < dirName {
			left = middle + 1
		} else {
			right = middle
		}

		if right-left < 1 {
			break
		}
	}

	return false
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

func (idx *Index) DeleteEntry(rootGoitPath string, path []byte) error {
	pos, _, isFound := idx.GetEntry(path)
	if !isFound {
		return fmt.Errorf("'%s' is not registered in index, so fail to delete", path)
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

func getEntriesFromTree(rootName string, nodes []*object.Node) ([]*Entry, error) {
	var entries []*Entry

	for _, node := range nodes {
		if len(node.Children) == 0 {
			var entryName string
			if rootName == "" {
				entryName = node.Name
			} else {
				entryName = fmt.Sprintf("%s/%s", rootName, node.Name)
			}
			newEntry := &Entry{
				Hash:       node.Hash,
				NameLength: uint16(len(entryName)),
				Path:       []byte(entryName),
			}
			entries = append(entries, newEntry)
		} else {
			var newRootName string
			if rootName == "" {
				newRootName = node.Name
			} else {
				newRootName = fmt.Sprintf("%s/%s", rootName, node.Name)
			}
			childEntries, err := getEntriesFromTree(newRootName, node.Children)
			if err != nil {
				return nil, err
			}
			entries = append(entries, childEntries...)
		}
	}

	return entries, nil
}

func (idx *Index) Reset(rootGoitPath string, hash sha.SHA1) error {
	// get commit
	commitObject, err := object.GetObject(rootGoitPath, hash)
	if err != nil {
		return fmt.Errorf("fail to get commit object: %w", err)
	}

	// get commit
	commit, err := object.NewCommit(commitObject)
	if err != nil {
		return fmt.Errorf("fail to get commit: %w", err)
	}

	// get tree object
	treeObject, err := object.GetObject(rootGoitPath, commit.Tree)
	if err != nil {
		return fmt.Errorf("fail to get tree object: %w", err)
	}

	// get tree
	tree, err := object.NewTree(rootGoitPath, treeObject)
	if err != nil {
		return fmt.Errorf("fail to get tree: %w", err)
	}

	// update index from tree
	entries, err := getEntriesFromTree("", tree.Children)
	if err != nil {
		return fmt.Errorf("fail to get entries from tree: %w", err)
	}
	idx.Header.EntryNum = uint32(len(entries))
	idx.Entries = entries

	// write index
	if err := idx.write(rootGoitPath); err != nil {
		return fmt.Errorf("fail to write index: %w", err)
	}

	return nil
}

func (idx *Index) DiffWithTree(tree *object.Tree) ([]*DiffEntry, error) {
	rootName := ""
	gotEntries, err := getEntriesFromTree(rootName, tree.Children)
	if err != nil {
		return nil, err
	}

	var diffEntries []*DiffEntry
	for _, gotEntry := range gotEntries {
		_, entry, isRegistered := idx.GetEntry(gotEntry.Path)
		if !isRegistered {
			diffEntries = append(diffEntries, &DiffEntry{
				Dt:    diffDelete,
				Entry: entry,
			})
		} else if !entry.Hash.Compare(gotEntry.Hash) {
			diffEntries = append(diffEntries, &DiffEntry{
				Dt:    diffModified,
				Entry: entry,
			})
		}
	}

	// check if there are new files
	for _, entry := range idx.Entries {
		_, isFound := object.GetNode(tree.Children, string(entry.Path))
		if !isFound {
			diffEntries = append(diffEntries, &DiffEntry{
				Dt:    diffNew,
				Entry: entry,
			})
		}
	}

	return diffEntries, nil
}
