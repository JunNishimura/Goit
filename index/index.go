package index

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/JunNishimura/Goit/sha"
)

type Entry struct {
	Hash       sha.SHA1
	NameLength uint16
	Path       string
}

func NewEntry(hash sha.SHA1, path string) *Entry {
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
	Entries []*Entry
}

func NewIndex() (*Index, error) {
	index := &Index{
		Header: Header{
			Signature: [4]byte{'D', 'I', 'R', 'C'},
			Version:   uint32(1),
			EntryNum:  uint32(0),
		},
	}
	if _, err := os.Stat(".goit/index"); !os.IsNotExist(err) {
		if err := index.read(); err != nil {
			return nil, fmt.Errorf("fail to read index: %v", err)
		}
	}
	return index, nil
}

func (idx *Index) IsUpdateNeeded() (bool, error) {
	return true, nil
}

func (idx *Index) Update(hash sha.SHA1, path string) error {
	entry := NewEntry(hash, path)
	idx.Entries = append(idx.Entries, entry)
	idx.EntryNum = uint32(len(idx.Entries))

	return idx.write()
}

func (idx *Index) read() error {
	return nil
}

func (idx *Index) write() error {
	f, err := os.Create(".goit/index")
	if err != nil {
		return fmt.Errorf("fail to create .goit/index: %v", err)
	}

	// fixed length encoding
	if err := binary.Write(f, binary.BigEndian, &idx.Header); err != nil {
		return fmt.Errorf("fail to write fixed-length encoding: %v", err)
	}

	// variable length encoding
	var data []byte
	for _, entry := range idx.Entries {
		bNameLength := make([]byte, 2)
		binary.BigEndian.PutUint16(bNameLength, entry.NameLength)
		data = append(data, entry.Hash...)
		data = append(data, bNameLength...)
		data = append(data, []byte(entry.Path)...)
	}
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("fail to write variable-length encoding: %v", err)
	}

	return nil
}
