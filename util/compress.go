package util

import (
	"bytes"
	"compress/zlib"
	"fmt"
)

func Compress(content string) (*bytes.Buffer, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer w.Close()
	if _, err := w.Write([]byte(content)); err != nil {
		return nil, fmt.Errorf("fail to write compressed data: %v", err)
	}
	return &b, nil
}
