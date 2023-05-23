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

	"github.com/JunNishimura/Goit/internal/binary"
	"github.com/JunNishimura/Goit/internal/sha"
)

type Object struct {
	Type Type
	Hash sha.SHA1
	Size int
	Data []byte
}

func NewObject(objType Type, data []byte) (*Object, error) {
	// get size of data
	size := len(data)

	// get hash of object
	checkSum := sha1.New()
	content := fmt.Sprintf("%s %d\x00%s", objType, size, data)
	_, err := io.WriteString(checkSum, content)
	if err != nil {
		return nil, err
	}
	hash := checkSum.Sum(nil)

	// make object
	object := &Object{
		Type: objType,
		Hash: hash,
		Size: size,
		Data: data,
	}

	return object, nil
}

func GetObject(rootGoitPath string, hash sha.SHA1) (*Object, error) {
	hashString := hash.String()
	objPath := filepath.Join(rootGoitPath, "objects", hashString[:2], hashString[2:])
	objFile, err := os.Open(objPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrIOHandling, objPath)
	}
	defer objFile.Close()

	zr, err := zlib.NewReader(objFile)
	if err != nil {
		return nil, fmt.Errorf("fail to construct zlib.NewReader: %w", err)
	}
	defer zr.Close()

	checkSum := sha1.New()
	tr := io.TeeReader(zr, checkSum)

	objType, size, err := readHeader(tr)
	if err != nil {
		return nil, fmt.Errorf("fail to read header: %w", err)
	}

	data, err := io.ReadAll(tr)
	if err != nil {
		return nil, ErrIOHandling
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
	// get header string
	headerStr, err := binary.ReadNullTerminatedString(r)
	if err != nil {
		return UndefinedObject, 0, ErrInvalidObject
	}

	// get type and size
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

func (o *Object) Write(rootGoitPath string) error {
	buf, err := o.compress()
	if err != nil {
		return err
	}

	dirPath := filepath.Join(rootGoitPath, "objects", o.Hash.String()[:2])
	filePath := filepath.Join(dirPath, o.Hash.String()[2:])
	if f, err := os.Stat(dirPath); os.IsNotExist(err) || !f.IsDir() {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("%w: %s", ErrIOHandling, dirPath)
		}
	}
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, filePath)
	}
	defer f.Close()
	if _, err := f.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("%w: %s", ErrIOHandling, filePath)
	}
	return nil
}
