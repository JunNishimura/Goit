package index

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/JunNishimura/Goit/sha"
)

type Entry struct {
	Hash sha.SHA1
	Path string
}

type Index struct {
	Entries []*Entry
}

type indexHeader struct {
	Signature [4]byte
	Version   uint32
	EntryNum  uint32
}

func (i *Index) Write() error {
	f, err := os.Create(".goit/index")
	if err != nil {
		return fmt.Errorf("fail to create .goit/index: %v", err)
	}

	// fixed length encoding
	header := indexHeader{
		Signature: [4]byte{'D', 'I', 'R', 'C'},
		Version:   uint32(1),
		EntryNum:  uint32(len(i.Entries)),
	}
	if err := binary.Write(f, binary.BigEndian, &header); err != nil {
		return fmt.Errorf("fail to write fixed-length encoding: %v", err)
	}

	// variable length encoding
	var data []byte
	for _, entry := range i.Entries {
		data = append(data, entry.Hash...)
		data = append(data, []byte(entry.Path)...)
		data = append(data, 0x00)
	}
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("fail to write variable-length encoding: %v", err)
	}

	return nil
}
