package index

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/JunNishimura/Goit/sha"
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

func (idx *Index) Update(hash sha.SHA1, path []byte) error {
	entry := NewEntry(hash, path)
	idx.Entries = append(idx.Entries, entry)
	idx.EntryNum = uint32(len(idx.Entries))

	return idx.write()
}

func (idx *Index) read() error {
	// read index
	b, err := ioutil.ReadFile(".goit/index")
	if err != nil {
		return fmt.Errorf("fail to read index: %v", err)
	}

	// make bytes reader
	buf := bytes.NewReader(b)

	// fixed length decoding
	err = binary.Read(buf, binary.BigEndian, &idx.Header)
	if err != nil {
		return fmt.Errorf("fail to read index header: %v", err)
	}

	// variable length decoding
	for i := 0; i < int(idx.EntryNum); i++ {
		// read hash
		hash := make(sha.SHA1, 20)
		err = binary.Read(buf, binary.BigEndian, &hash)
		if err != nil {
			return fmt.Errorf("fail to read hash from index: %v", err)
		}

		// read file name length
		var nameLength uint16
		err = binary.Read(buf, binary.BigEndian, &nameLength)
		if err != nil {
			return fmt.Errorf("fail to read file name length from index: %v", err)
		}

		// read file path
		path := make([]byte, nameLength)
		err = binary.Read(buf, binary.BigEndian, &path)
		if err != nil {
			return fmt.Errorf("fail to read path from index: %v", err)
		}

		entry := NewEntry(hash, path)
		idx.Entries = append(idx.Entries, entry)
	}

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
		data = append(data, entry.Path...)
	}
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("fail to write variable-length encoding: %v", err)
	}

	return nil
}
