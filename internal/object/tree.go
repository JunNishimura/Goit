package object

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/JunNishimura/Goit/internal/binary"
	"github.com/JunNishimura/Goit/internal/sha"
	index "github.com/JunNishimura/Goit/internal/store"
)

func WriteTreeObject(rootGoitPath string, entries []*index.Entry) (*Object, error) {
	var dirName string
	var data []byte
	var entryBuf []*index.Entry
	i := 0
	for {
		if i >= len(entries) {
			// if the last entry is in the directory
			if dirName != "" {
				treeObject, err := WriteTreeObject(rootGoitPath, entryBuf)
				if err != nil {
					return nil, err
				}
				data = append(data, []byte(fmt.Sprintf("040000 %s", dirName))...)
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
				treeObject, err := WriteTreeObject(rootGoitPath, entryBuf)
				if err != nil {
					return nil, err
				}
				data = append(data, []byte(fmt.Sprintf("040000 %s", dirName))...)
				data = append(data, 0x00)
				data = append(data, treeObject.Hash...)
				// clear dirName and entryBuf
				dirName = ""
				entryBuf = make([]*index.Entry, 0)
			} else {
				data = append(data, []byte(fmt.Sprintf("100644 %s", string(entry.Path)))...)
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
				treeObject, err := WriteTreeObject(rootGoitPath, entryBuf)
				if err != nil {
					return nil, err
				}
				data = append(data, []byte(fmt.Sprintf("040000 %s", dirName))...)
				data = append(data, 0x00)
				data = append(data, treeObject.Hash...)
				// clear dirName and entryBuf
				dirName = ""
				entryBuf = make([]*index.Entry, 0)
			}
		}
	}

	// make tree object
	treeObject, err := NewObject(TreeObject, data)
	if err != nil {
		return nil, err
	}

	// write tree object
	if err := treeObject.Write(rootGoitPath); err != nil {
		return nil, err
	}

	return treeObject, nil
}

func (to *Object) ExtractEntries(rootGoitPath, rootDir string) ([]*index.Entry, error) {
	var entries []*index.Entry
	var dirName string
	var filePath string
	isFirstLine := true

	buf := bytes.NewReader(to.Data)
	for {
		var lineSplit []string
		if isFirstLine {
			lineString, err := binary.ReadNullTerminatedString(buf)
			if err != nil {
				return nil, err
			}
			lineSplit = strings.Split(lineString, " ")
			isFirstLine = false
		} else {
			// read 20 bytes sha1 hash
			hashBytes := make([]byte, 20)
			n, err := buf.Read(hashBytes)
			if err != nil {
				return nil, err
			}
			if n != 20 {
				return nil, errors.New("fail to read hash")
			}

			// read filemode and path
			lineString, err := binary.ReadNullTerminatedString(buf)
			if err != nil {
				return nil, err
			}

			lineSplit = []string{string(hashBytes)}
			if lineString != "" {
				lineSplit = append(lineSplit, strings.Split(lineString, " ")...)
			}
		}

		if dirName != "" {
			hashString := hex.EncodeToString([]byte(lineSplit[0]))
			hash, err := sha.ReadHash(hashString)
			if err != nil {
				return nil, err
			}
			treeObject, err := GetObject(rootGoitPath, hash)
			if err != nil {
				return nil, err
			}
			var path string
			if rootDir == "" {
				path = dirName
			} else {
				path = fmt.Sprintf("%s/%s", rootDir, dirName)
			}
			getEntries, err := treeObject.ExtractEntries(rootGoitPath, path)
			if err != nil {
				return nil, err
			}
			entries = append(entries, getEntries...)
			dirName = ""
		}

		if filePath != "" {
			hashString := hex.EncodeToString([]byte(lineSplit[0]))
			hash, err := sha.ReadHash(hashString)
			if err != nil {
				return nil, err
			}
			entry := index.NewEntry(hash, []byte(filePath))
			entries = append(entries, entry)
			filePath = ""
		}

		if len(lineSplit) == 1 { // last line
			break
		} else {
			var fileMode string
			var fileName string
			if len(lineSplit) == 2 {
				fileMode = lineSplit[0]
				fileName = lineSplit[1]
			} else if len(lineSplit) == 3 {
				fileMode = lineSplit[1]
				fileName = lineSplit[2]
			}

			if fileMode == "040000" {
				dirName = fileName
			} else if fileMode == "100644" {
				if rootDir == "" {
					filePath = fileName
				} else {
					filePath = fmt.Sprintf("%s/%s", rootDir, fileName)
				}
			}
		}
	}
	return entries, nil
}

func (to *Object) ConvertDataToString() (string, error) {
	var lines []string
	isFirstLine := true

	buf := bytes.NewReader(to.Data)
	for {
		if isFirstLine {
			lineString, err := binary.ReadNullTerminatedString(buf)
			if err != nil {
				return "", err
			}
			if lineString == "" {
				return "", nil
			}
			lines = append(lines, lineString)
			isFirstLine = false
		} else {
			// read hash
			hashBytes := make([]byte, 20)
			n, err := buf.Read(hashBytes)
			if err != nil {
				return "", err
			}
			if n != 20 {
				return "", errors.New("fail to read hash bytes")
			}
			hashString := hex.EncodeToString(hashBytes)
			lines = append(lines, hashString)

			// read filemode and path
			lineString, err := binary.ReadNullTerminatedString(buf)
			if err != nil {
				return "", err
			}
			if lineString == "" {
				break
			} else {
				lines = append(lines, lineString)
			}
		}
	}

	dataString := strings.Join(lines, "\n")

	return dataString, nil
}
