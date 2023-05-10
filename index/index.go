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

type Header struct {
	Signature [4]byte
	Version   uint32
	EntryNum  uint32
}

type Index struct {
	Header  Header
	Entries []*Entry
}

func (i *Index) Write() error {
	f, err := os.Create(".goit/index")
	if err != nil {
		return fmt.Errorf("fail to create .goit/index: %v", err)
	}

	// fixed length encoding
	if err := binary.Write(f, binary.BigEndian, &i.Header); err != nil {
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
