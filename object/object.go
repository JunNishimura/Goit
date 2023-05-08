package object

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"path/filepath"

	"github.com/JunNishimura/Goit/sha"
)

type Object struct {
	Type Type
	Hash sha.SHA1
	Size int
	Data []byte
}

func (o *Object) Header() []byte {
	return []byte(fmt.Sprintf("%s %d\x00", o.Type, o.Size))
}

func (o *Object) Write(compData []byte) error {
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
	if _, err := f.Write(compData); err != nil {
		return fmt.Errorf("fail to write to %s: %v", filePath, err)
	}
	return nil
}

func (o *Object) CompressBlob() ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer w.Close()
	data := append(o.Header(), o.Data...)
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("fail to compress data: %v", err)
	}
	return b.Bytes(), nil
}
