package object

import (
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/JunNishimura/Goit/index"
)

func NewTreeObject(indexClient *index.Index) *Object {
	var data []byte
	for _, entry := range indexClient.Entries {
		data = append(data, entry.Path...)
		data = append(data, 0x00)
		data = append(data, entry.Hash...)
	}

	// get size of data
	size := len(data)

	// get hash value
	checkSum := sha1.New()
	content := fmt.Sprintf("%s %d\x00%s", TreeObject, size, data)
	io.WriteString(checkSum, content)
	hash := checkSum.Sum(nil)

	object := &Object{
		Type: TreeObject,
		Hash: hash,
		Size: size,
		Data: data,
	}

	return object
}
