package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/JunNishimura/Goit/index"
	"github.com/JunNishimura/Goit/sha"
)

type Object struct {
	Type Type
	Hash sha.SHA1
	Size int
	Data []byte
}

func NewObject(objType Type, data []byte) *Object {
	// get size of data
	size := len(data)

	// get hash of object
	checkSum := sha1.New()
	content := fmt.Sprintf("%s %d\x00%s", objType, size, data)
	io.WriteString(checkSum, content)
	hash := checkSum.Sum(nil)

	// make object
	object := &Object{
		Type: objType,
		Hash: hash,
		Size: size,
		Data: data,
	}

	return object
}

func (o *Object) Header() []byte {
	return []byte(fmt.Sprintf("%s %d\x00", o.Type, o.Size))
}

func (o *Object) compress() (bytes.Buffer, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	data := append(o.Header(), o.Data...)
	if _, err := w.Write(data); err != nil {
		return b, fmt.Errorf("fail to compress data: %v", err)
	}
	w.Close()
	return b, nil
}

func (o *Object) Write() error {
	buf, err := o.compress()
	if err != nil {
		return err
	}

	dirPath := filepath.Join(".goit", "objects", o.Hash.String()[:2])
	filePath := filepath.Join(dirPath, o.Hash.String()[2:])
	if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("fail to make %s: %v", dirPath, err)
	}
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("fail to make %s: %v", filePath, err)
	}
	defer f.Close()
	if _, err := f.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("fail to write to %s: %v", filePath, err)
	}
	return nil
}

func MakeTreeObject(entries []*index.Entry) *Object {
	var dirName string
	var data []byte
	var entryBuf []*index.Entry
	i := 0
	for {
		if i >= len(entries) {
			// if the last entry is in the directory
			if dirName != "" {
				treeObject := MakeTreeObject(entryBuf)
				data = append(data, []byte(dirName)...)
				data = append(data, 0x00)
				data = append(data, treeObject.Hash...)
			}
			break
		}

		entry := entries[i]
		slashSplit := strings.SplitN(string(entry.Path), "/", 2)
		if len(slashSplit) == 1 {
			if dirName != "" {
				// make tree object from entryBuf
				treeObject := MakeTreeObject(entryBuf)
				data = append(data, []byte(dirName)...)
				data = append(data, 0x00)
				data = append(data, treeObject.Hash...)
				// clear dirName and entryBuf
				dirName = ""
				entryBuf = make([]*index.Entry, 0)
			} else {
				data = append(data, entry.Path...)
				data = append(data, 0x00)
				data = append(data, entry.Hash...)
				i++
			}
		} else {
			if dirName == "" {
				dirName = slashSplit[0]
				newEntry := index.NewEntry(entry.Hash, []byte(slashSplit[1]))
				entryBuf = append(entryBuf, newEntry)
				i++
			} else if dirName != "" && dirName == slashSplit[0] {
				// same dir with prev entry
				newEntry := index.NewEntry(entry.Hash, []byte(slashSplit[1]))
				entryBuf = append(entryBuf, newEntry)
				i++
			} else if dirName != "" && dirName != slashSplit[0] {
				treeObject := MakeTreeObject(entryBuf)
				data = append(data, []byte(dirName)...)
				data = append(data, 0x00)
				data = append(data, treeObject.Hash...)
				// clear dirName and entryBuf
				dirName = ""
				entryBuf = make([]*index.Entry, 0)
			}
		}
	}

	// make tree object
	treeObject := NewObject(TreeObject, data)

	return treeObject
}

func (o *Object) UpdateBranch() error {
	filePath := filepath.Join(".goit", "refs", "heads", "main")
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("fail to make %s: %v", filePath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(o.Hash.String()); err != nil {
		return fmt.Errorf("fail to write hash to %s: %v", filePath, err)
	}

	return nil
}
