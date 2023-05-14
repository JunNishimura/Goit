package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
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

func GetObject(hash sha.SHA1) (*Object, error) {
	hashString := hash.String()
	objPath := filepath.Join(".goit", "objects", hashString[:2], hashString[2:])
	objFile, err := os.Open(objPath)
	if err != nil {
		return nil, fmt.Errorf("fail to open %s: %v", objPath, err)
	}
	defer objFile.Close()

	zr, err := zlib.NewReader(objFile)
	if err != nil {
		return nil, fmt.Errorf("fail to construct zlib.NewReader: %v", err)
	}
	defer zr.Close()

	checkSum := sha1.New()
	tr := io.TeeReader(zr, checkSum)

	objType, size, err := readHeader(tr)
	if err != nil {
		return nil, fmt.Errorf("fail to read header: %v", err)
	}

	data, err := ioutil.ReadAll(tr)
	if err != nil {
		return nil, fmt.Errorf("fail to read object: %v", err)
	}

	if len(data) != size {
		return nil, ErrInvalidObject
	}

	objHash := checkSum.Sum(nil)

	object := &Object{
		Type: objType,
		Hash: objHash,
		Size: size,
		Data: data,
	}

	return object, nil
}

func readHeader(r io.Reader) (Type, int, error) {
	// read until null byte
	headerBytes := make([]byte, 0)
	for {
		b := make([]byte, 1)
		_, err := r.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			return UndefinedObject, 0, err
		}
		if b[0] == 0 {
			break
		}
		headerBytes = append(headerBytes, b[0])
	}

	// get type and size
	headerStr := string(headerBytes)
	headerSplit := strings.SplitN(headerStr, " ", 2)
	if len(headerSplit) != 2 {
		return UndefinedObject, 0, ErrInvalidObject
	}

	objTypeString := headerSplit[0]
	sizeString := headerSplit[1]

	objType, err := NewType(objTypeString)
	if err != nil {
		return UndefinedObject, 0, err
	}
	var size int
	if _, err := fmt.Sscanf(sizeString, "%d", &size); err != nil {
		return UndefinedObject, 0, err
	}

	return objType, size, nil
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
