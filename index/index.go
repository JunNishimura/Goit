package index

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/JunNishimura/Goit/sha"
)

type Entry struct {
	NameLength uint16
	Hash       sha.SHA1
	Path       string
}

type Header struct {
	Signature [4]byte
	Version   uint32
	EntryNum  uint32
}

type Index struct {
	Header  Header
	Entries []*Entry
}

func NewIndex() *Index {
	return &Index{
		Header: Header{
			Signature: [4]byte{'D', 'I', 'R', 'C'},
			Version:   uint32(1),
			EntryNum:  uint32(0),
		},
	}
}

func (idx *Index) IsUpdateNeeded() (bool, error) {
	return false, nil
}

func (idx *Index) Update() error {
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
		data = append(data, entry.Hash...)
		data = append(data, []byte(entry.Path)...)
	}
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("fail to write variable-length encoding: %v", err)
	}

	return nil
}
