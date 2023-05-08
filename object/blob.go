package object

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func NewBlobObject(filePath string) (*Object, error) {
	f, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf(`fatal: Cannot open '%s': No such file`, filePath)
	}
	if f.IsDir() {
		return nil, fmt.Errorf(`fatal: '%s' is invalid to make blob object`, filePath)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("fail to read file: %v", err)
	}
	size := len(data)

	checkSum := sha1.New()
	content := fmt.Sprintf("%s %d\x00%s", BlobObject, size, data)
	io.WriteString(checkSum, content)
	hash := checkSum.Sum(nil)

	object := &Object{
		Type: BlobObject,
		Hash: hash,
		Size: size,
		Data: data,
	}

	return object, nil
}
